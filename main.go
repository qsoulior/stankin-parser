package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/qsoulior/stankin-parser/schedule"
	pdf_decoder "github.com/qsoulior/stankin-parser/schedule/decoder/pdf"
	"github.com/qsoulior/stankin-parser/schedule/encoder"
	ical_encoder "github.com/qsoulior/stankin-parser/schedule/encoder/ical"
	json_encoder "github.com/qsoulior/stankin-parser/schedule/encoder/json"
)

func decode() ([]schedule.Unit, *schedule.Meta, error) {
	reader, err := os.Open(input)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = reader.Close() }()

	stat, err := reader.Stat()
	if err != nil {
		return nil, nil, err
	}

	decoder := pdf_decoder.New(reader, stat.Size())
	units, meta, err := decoder.Decode()
	if err != nil {
		return nil, nil, err
	}

	return units, meta, nil
}

func getEncoder(format string, w io.Writer) (encoder.Encoder, error) {
	switch format {
	case "ical":
		return ical_encoder.New(w), nil
	case "json":
		return json_encoder.New(w), nil
	default:
		return nil, errors.New("invalid encoding format")
	}
}

func encode(events []schedule.Event, group string) error {
	writer, err := os.Create(output)
	if err != nil {
		return err
	}
	defer func() { _ = writer.Close() }()

	encoder, err := getEncoder(format, writer)
	if err != nil {
		return err
	}

	return encoder.Encode(events, group, schedule.EventSubgroup(subgroup))
}

var (
	input, output string
	format        string
	subgroup      string
	year          int
)

func init() {
	// required
	flag.StringVar(&input, "input", "", "input file path")
	flag.StringVar(&output, "output", "", "output file path")

	// optional
	flag.StringVar(&format, "format", "ical", "encoding format")
	flag.StringVar(&subgroup, "subgroup", "", "subgroup to filter events")
	flag.IntVar(&year, "year", time.Now().Year(), "year of events")
}

func main() {
	log := log.New(os.Stdout, "", 0)

	flag.Parse()

	if input == "" || output == "" {
		flag.PrintDefaults()
		log.Fatal("expected input and output")
	}

	// decode
	units, meta, err := decode()
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}

	log.Printf("schedule decoded from %s\n", input)
	log.Printf("name of decoded group: %s\n", meta.Group)
	log.Printf("number of decoded units: %d\n", len(units))

	// parse
	events, err := schedule.Parse(units, year)
	if err != nil {
		log.Fatalf("failed to parse: %s", err)
	}

	log.Printf("number of parsed events: %d\n", len(events))

	// encode
	err = encode(events, meta.Group)
	if err != nil {
		log.Fatalf("failed to encode: %s", err)
	}

	log.Printf("schedule encoded into %s\n", output)
}
