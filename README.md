# stankin-parser
[![Go Version](https://img.shields.io/github/go-mod/go-version/qsoulior/stankin-parser?style=flat-square)](https://go.dev/doc/go1.22)
[![Go Report Card](https://goreportcard.com/badge/github.com/qsoulior/stankin-parser?style=flat-square)](https://goreportcard.com/report/github.com/qsoulior/stankin-parser)

stankin-parser is CLI utility to parse stankin schedules and convert them to json or ical formats.

## Running
```
go run . -input foo.pdf -output bar.ics -format ical
```

## Flags
| Flag | Default value | Description |
| ---- | ------------- | ----------- |
| input || input file path |
| output || output file path |
| format | ical | encoding format (must be "ical" or "json") |
| subgroup | all subgroups | subgroup to filter events |
| year | current year | year of events since source schedule doesn't contain it |
