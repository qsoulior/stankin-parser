package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/qsoulior/stankin-parser/schedule"
	"github.com/qsoulior/stankin-parser/schedule/decoder"
	"github.com/qsoulior/stankin-parser/schedule/encoder"
)

func decode() ([]schedule.Cell, *schedule.Meta, error) {
	reader, err := os.Open(input)
	if err != nil {
		return nil, nil, err
	}
	defer reader.Close()

	stat, err := reader.Stat()
	if err != nil {
		log.Fatal(err)
	}

	decoder := decoder.NewPDF(reader, stat.Size())
	cells, meta, err := decoder.Decode()
	if err != nil {
		return nil, nil, err
	}

	return cells, meta, nil
}

func getEncoder(format string, w io.Writer) (encoder.Encoder, error) {
	switch format {
	case "ical":
		return encoder.NewIcal(w), nil
	case "json":
		return encoder.NewJSON(w), nil
	default:
		return nil, errors.New("invalid encoding format")
	}
}

func encode(events []schedule.Event, group string) error {
	writer, err := os.Create(output)
	if err != nil {
		return err
	}
	defer writer.Close()

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
	cells, meta, err := decode()
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}

	log.Printf("schedule decoded from %s\n", input)
	log.Printf("name of decoded group: %s\n", meta.Group)
	log.Printf("number of decoded cells: %d\n", len(cells))

	// parse
	events, err := schedule.Parse(cells, year)
	if err != nil {
		log.Fatalf("failed to parse: %s", err)
	}

	log.Printf("number of parsed events: %d\n", len(events))

	// encode
	err = encode(events, meta.Group)
	if err != nil {
		log.Fatalf("failed to encode: %s", err)
	}

	log.Printf("schedule encoded to %s\n", output)
}
