package main

import "archive/zip"
import "encoding/xml"
import "fmt"
import "io"
import "log"
import "os"
import "regexp"
import "sort"
import "time"
import "golang.org/x/text/encoding/unicode"
import "golang.org/x/text/transform"

// Parse options
import "github.com/DavidGamba/go-getoptions"

// GORM
import "gorm.io/driver/sqlite"
import "gorm.io/gorm"
import "gorm.io/gorm/clause"
import "gorm.io/gorm/logger"

var clock time.Time

func main() {
	var skipsolar bool
	var skipbattery bool
	var zipfile string
	var outfile string
	var batchsize int
	opt := getoptions.New()
	opt.StringVar(&zipfile, "z", "", opt.Alias("Z"), opt.Description("Name der Export-ZIP-Datei"))
	opt.StringVar(&outfile, "o", "mastr.sql", opt.Alias("O"), opt.Description("Ausgabedatei (SQLite)"))
	opt.BoolVar(&skipsolar, "S", false, opt.Description("Solardaten nicht laden"))
	opt.BoolVar(&skipbattery, "A", false, opt.Description("Speicherdaten nicht laden"))
	opt.IntVar(&batchsize, "b", 500, opt.Description("Batch Insert Count"))
	remaining, _ := opt.Parse(os.Args[1:])

	if len(zipfile) == 0 {
		fmt.Println("ERROR: No ZIP file provided")
		fmt.Println()
		fmt.Println(opt.Help())
		os.Exit(1)
	}

	if len(remaining) > 0 {
		fmt.Println(opt.Help())
		os.Exit(1)
	}

	esfiles_re := regexp.MustCompile(`EinheitenSolar_\d+.xml$`)
	ssfiles_re := regexp.MustCompile(`EinheitenStromSpeicher_\d+.xml$`)

	customlogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{SlowThreshold: 5 * time.Second})

	db, err := gorm.Open(sqlite.Open(outfile), &gorm.Config{Logger: customlogger, CreateBatchSize: batchsize, SkipDefaultTransaction: true})
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("PRAGMA journal_mode = WAL;")
	db.Exec("PRAGMA synchronous = NORMAL;")
	db.AutoMigrate(&EinheitSolar{})
	db.AutoMigrate(&EinheitStromSpeicher{})

	zr, err := zip.OpenReader(zipfile)
	if err != nil {
		log.Fatal(err)
	}
	defer zr.Close()

	sort.Slice(zr.File, func(i, j int) bool {
		return zr.File[i].Name < zr.File[j].Name
	})

	for _, file := range zr.File {
		filename := file.Name

		if esfiles_re.MatchString(filename) {
			if skipsolar {
				continue
			}

			clock = time.Now()
			fmt.Printf("Processing %s\n", filename)

			content, err := GetZipFileContent(file)
			if err != nil {
				log.Printf("Could not read %s: %s\n", filename, err.Error())
				continue
			}

			es, err := ParseEinheitenSolar(content)
			if err != nil {
				log.Printf("Could not parse %s: %s\n", filename, err.Error())
				content.Close()
				continue
			}

			fmt.Printf("Parsed out %d records\n", len(es))

			// Geändert: OnConflict fängt doppelte IDs aus verschiedenen XML-Dateien ab
			result := db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(es, batchsize)
			content.Close()
			
			if result.Error != nil {
				log.Fatalf("DB error: %s", result.Error)
			}

			rows := result.RowsAffected

			fmt.Printf("Processing time: %s\n", time.Since(clock))

			if int64(len(es)) != rows {
				fmt.Printf("Not all data could be loaded: %d units, %d rows affected\n", len(es), rows)
			}
		}

		// Speicher
		if ssfiles_re.MatchString(filename) {
			if skipbattery {
				continue
			}

			clock = time.Now()
			fmt.Printf("Processing %s\n", filename)

			content, err := GetZipFileContent(file)
			if err != nil {
				log.Printf("Could not read %s: %s\n", filename, err.Error())
				continue
			}

			ss, err := ParseStromSpeicher(content)
			if err != nil {
				log.Printf("Could not parse %s: %s\n", filename, err.Error())
				content.Close()
				continue
			}

			fmt.Printf("Parsed out %d records\n", len(ss))

			// Geändert: OnConflict fängt doppelte IDs aus verschiedenen XML-Dateien ab
			result := db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(ss, batchsize)
			content.Close()
			
			if result.Error != nil {
				log.Fatalf("DB error: %s", result.Error)
			}

			rows := result.RowsAffected

			fmt.Printf("Processing time: %s\n", time.Since(clock))

			if int64(len(ss)) != rows {
				fmt.Printf("Not all data could be loaded: %d units, %d rows affected\n", len(ss), rows)
			}
		}
	}

	tx := db.Begin()
	tx.Exec(`UPDATE einheit_solars SET bundesland = "BB" WHERE bundesland = "1400"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "BE" WHERE bundesland = "1401"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "BW" WHERE bundesland = "1402"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "BY" WHERE bundesland = "1403"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "HE" WHERE bundesland = "1404"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "HB" WHERE bundesland = "1405"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "HH" WHERE bundesland = "1406"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "MV" WHERE bundesland = "1407"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "NI" WHERE bundesland = "1408"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "NW" WHERE bundesland = "1409"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "RP" WHERE bundesland = "1410"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "SH" WHERE bundesland = "1411"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "SL" WHERE bundesland = "1412"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "SN" WHERE bundesland = "1413"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "ST" WHERE bundesland = "1414"`)
	tx.Exec(`UPDATE einheit_solars SET bundesland = "TH" WHERE bundesland = "1415"`)

	tx.Exec(`UPDATE einheit_solars SET land = "DE" WHERE land = "84"`)
	tx.Exec(`UPDATE einheit_solars SET land = "AT" WHERE land = "206"`)
	tx.Exec(`UPDATE einheit_solars SET land = "CH" WHERE land = "231"`)
	tx.Exec(`UPDATE einheit_solars SET land = "DK" WHERE land = "90"`)
	tx.Exec(`UPDATE einheit_solars SET land = "BE" WHERE land = "66"`)
	tx.Commit()
}

func GetZipFileContent(file *zip.File) (io.ReadCloser, error) {
	return file.Open()
}

func ParseEinheitenSolar(xmlStream io.Reader) ([]*EinheitSolar, error) {
	decoder := xml.NewDecoder(TransformUTF16(xmlStream))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return input, nil
	}

	var einheiten []*EinheitSolar
	clock = time.Now()

	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Could not parse XML token: %s\n", err.Error())
			return nil, err
		}

		// Sucht gezielt nach dem Start-Element eines Eintrags
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "EinheitSolar" {
			var eintrag EinheitSolar
			if err := decoder.DecodeElement(&eintrag, &se); err != nil {
				return nil, err
			}
			einheiten = append(einheiten, &eintrag)
		}
	}

	fmt.Printf("Decoding took %s\n", time.Since(clock))
	return einheiten, nil
}

func ParseStromSpeicher(xmlStream io.Reader) ([]*EinheitStromSpeicher, error) {
	decoder := xml.NewDecoder(TransformUTF16(xmlStream))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return input, nil
	}

	var einheiten []*EinheitStromSpeicher
	clock = time.Now()

	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Could not parse XML token: %s\n", err.Error())
			return nil, err
		}

		// Sucht gezielt nach dem Start-Element eines Eintrags
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "EinheitStromSpeicher" {
			var eintrag EinheitStromSpeicher
			if err := decoder.DecodeElement(&eintrag, &se); err != nil {
				return nil, err
			}
			einheiten = append(einheiten, &eintrag)
		}
	}

	fmt.Printf("Decoding took %s\n", time.Since(clock))
	return einheiten, nil
}

func TransformUTF16(xmlStream io.Reader) io.Reader {
	utf16bom := unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
	transformer := utf16bom.NewDecoder()
	return transform.NewReader(xmlStream, transformer)
}
