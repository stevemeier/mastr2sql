# mastr2sql
Transform Marktstammdatenregister (MaStR) data into an SQLite database

## Why?

With the surge of solar installtions in Germany, it has been a topic of many discussions.
As I was curious about this topic, I wanted to dig deeper into the public data.
As it is published in the most inconvenient format possible (UTF-16 encoded XML), I decided to write this tool to turn it into a database.

## Terms of use

This piece of software comes with no guarantees whatsoever. Use at your own risk.
If you identify any major issues or if you have any cool expansion ideas, please create a GitHub issue.
It is published under the MIT license, so you can pretty much use it as you like.

## How to use

### Build the binary

`go build -o mastr2sql *.go`

### Download public data snapshot

Download the most recent ZIP file from this URL:
https://www.marktstammdatenregister.de/MaStR/Datendownload (around 3 GB)

### Run an import

`mastr2sql -z Gesamtdatenexport_XYZ.zip`

## Important notes

1. Data is not sanity-checked, it is merely imported
2. Fields `bundesland` and `land` are transformed for ease of use.
  After the import, these will have ISO 3166-2 and ISO 3166-2:DE values instead its numeric original

## Command line options

```
SYNOPSIS:
    mastr2sql [-A] [-S] [-b <int>] [-o|-O <string>] [-z|-Z <string>]
                       [<args>]

OPTIONS:
    -A                Do NOT load data on batteries (default: false)

    -S                Do NOT load data on solar installations (default: false)

    -b <int>          Batch Insert Count (default: 500)

    -o|-O <string>    Output file (SQLite) (default: "mastr.sql")

    -z|-Z <string>    ZIP file to load data from (default: "")
```

## Database structure

```
CREATE TABLE `einheit_solars` (
  `einheit_mastr_nummer` text,
  `datum_letzte_aktualisierung` text,
  `lokation_mastr_nummer` text,
  `netzbetreiberpruefung_status` integer,
  `netzbetreiberpruefung_datum` text,
  `anlagenbetreiber_mastr_nummer` text,
  `land` integer,
  `bundesland` integer,
  `landkreis` text,
  `gemeinde` text,
  `gemeindeschluessel` text,
  `postleitzahl` text,
  `strasse` text,
  `gemarkung` text,
  `flur_flurstuecknummern` text,
  `hausnummer` text,
  `adresszusatz` text,
  `ort` text,
  `laengengrad` real,
  `breitengrad` real,
  `registrierungsdatum` text,
  `inbetriebnahmedatum` text,
  `datum_endgueltige_stilllegung` text,
  `datum_beginn_voruebergehende_stilllegung` text,
  `datum_wiederaufnahme_betrieb` text,
  `geplantes_inbetriebnahmedatum` text,
  `einheit_systemstatus` integer,
  `einheit_betriebsstatus` integer,
  `bestandsanlage_mastr_nummer` text,
  `alt_anlagenbetreiber_mastr_nummer` text,
  `datum_des_betreiberwechsels` text,
  `datum_registrierung_des_betreiberwechsels` text,
  `name_stromerzeugungseinheit` text,
  `weic` text,
  `weicnv` integer,
  `weic_display_name` text,
  `kraftwerksnummer` text,
  `kraftwerksnummernv` integer,
  `energietraeger` integer,
  `bruttoleistung` real,
  `nettonennleistung` real,
  `anschluss_an_hoechst_oder_hoch_spannung` integer,
  `fernsteuerbarkeit_nb` integer,
  `fernsteuerbarkeit_dv` integer,
  `einspeisungsart` integer,
  `zugeordnete_wirkleistung_wechselrichter` real,
  `anzahl_module` integer,
  `lage` integer,
  `leistungsbegrenzung` integer,
  `einheitliche_ausrichtung_und_neigungswinkel` integer,
  `hauptausrichtung` integer,
  `hauptausrichtung_neigungswinkel` integer,
  `nebenausrichtung` integer,
  `nebenausrichtung_neigungswinkel` integer,
  `nutzungsbereich` integer,
  `buergerenergie` integer,
  `eeg_mastr_nummer` text,
  `in_anspruch_genommene_flaeche` real,
  `art_der_flaeche_ids` text,
  `in_anspruch_genommene_ackerflaeche` real,
  `einsatzverantwortlicher` text,
  `gen_mastr_nummer` text,

  PRIMARY KEY (`einheit_mastr_nummer`)
);

CREATE TABLE `einheit_strom_speichers` (
  `einheit_mastr_nummer` text,
  `datum_letzte_aktualisierung` text,
  `lokation_ma_st_r_nummer` text,
  `netzbetreiberpruefung_status` integer,
  `netzbetreiberpruefung_datum` text,
  `anlagenbetreiber_mastr_nummer` text,
  `land` integer,
  `bundesland` integer,
  `landkreis` text,
  `gemeinde` text,
  `gemeindeschluessel` text,
  `postleitzahl` text,
  `gemarkung` text,
  `flur_flurstuecknummern` text,
  `strasse` text,
  `hausnummer` text,
  `adresszusatz` text,
  `ort` text,
  `laengengrad` real,
  `breitengrad` real,
  `registrierungsdatum` text,
  `inbetriebnahmedatum` text,
  `geplantes_inbetriebnahmedatum` text,
  `datum_endgueltige_stilllegung` text,
  `datum_beginn_voruebergehende_stilllegung` text,
  `datum_wiederaufnahme_betrieb` text,
  `einheit_systemstatus` integer,
  `einheit_betriebsstatus` integer,
  `bestandsanlage_mastr_nummer` text,
  `alt_anlagenbetreiber_mastr_nummer` text,
  `datum_des_betreiberwechsels` text,
  `datum_registrierung_des_betreiberwechsels` text,
  `name_stromerzeugungseinheit` text,
  `weic` text,
  `weicnv` integer,
  `weic_display_name` text,
  `kraftwerksnummer` text,
  `kraftwerksnummernv` integer,
  `energietraeger` integer,
  `bruttoleistung` real,
  `nettonennleistung` real,
  `anschluss_an_hoechst_oder_hoch_spannung` integer,
  `einsatzverantwortlicher` text,
  `fernsteuerbarkeit_nb` integer,
  `fernsteuerbarkeit_dv` integer,
  `einspeisungsart` integer,
  `einsatzort` integer,
  `gen_mastr_nummer` text,
  `reserveart_nach_dem_en_wg` integer,
  `datum_ueberfuehrung_in_reserve` text,
  `ac_dc_koppelung` integer,
  `batterietechnologie` integer,
  `notstromaggregat` integer,
  `pumpbetrieb_leistungsaufnahme` real,
  `pumpbetrieb_kontinuierlich_regelbar` integer,
  `pumpspeichertechnologie` integer,
  `bestandteil_grenzkraftwerk` integer,
  `nettonennleistung_deutschland` real,
  `zugeordnente_wirkleistung_wechselrichter` real,
  `spe_mastr_nummer` text,
  `eeg_ma_st_r_nummer` text,
  `eeg_anlagentyp` integer,
  `technologie` integer,
  `gemeinsam_registrierte_solareinheit_mastr_nummer` text,

  PRIMARY KEY (`einheit_mastr_nummer`)
);
```
