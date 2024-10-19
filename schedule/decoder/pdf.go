package decoder

import (
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/qsoulior/stankin-parser/schedule"
)

const pdfPageNum = 1

type pdfDecoder struct {
	r    io.ReaderAt
	size int64
}

// NewPDF creates and returns new pdf decoder.
func NewPDF(r io.ReaderAt, size int64) *pdfDecoder {
	return &pdfDecoder{r, size}
}

func (d *pdfDecoder) decodeMeta(chunks []pdf.Text) (*schedule.Meta, int) {
	i := 1
	for chunks[i].Y == chunks[i-1].Y && i < len(chunks) {
		i++
	}

	var data strings.Builder
	data.Grow(i - 1)
	for j := 0; j < i; j++ {
		data.WriteString(chunks[j].S)
	}

	meta := schedule.Meta{
		Group: data.String(),
	}

	for chunks[i].X < 42 || chunks[i].Y > 520 && i < len(chunks) {
		i++
	}

	return &meta, i
}

func (d *pdfDecoder) decodeUnit(chunks []pdf.Text) (schedule.Unit, int) {
	var data strings.Builder
	data.Grow(len(chunks[0].S))
	data.WriteString(chunks[0].S)

	i := 1
	iMax := 0
	for chunks[i-1].S != "]" && i < len(chunks) {
		if chunks[i].X > chunks[iMax].X {
			iMax = i
		}

		if chunks[i].Y != chunks[i-1].Y {
			data.WriteRune(' ')
		}
		data.WriteString(chunks[i].S)
		i++
	}

	unit := schedule.Unit{
		Data:   data.String(),
		Left:   int(chunks[0].X),
		Top:    int(chunks[0].Y),
		Right:  int(chunks[iMax].X),
		Bottom: int(chunks[i-1].Y),
	}

	return unit, i
}

// Decode decodes schedule units and metadata and returns them.
// It uses [io.ReaderAt] as source of input data.
func (d *pdfDecoder) Decode() ([]schedule.Unit, *schedule.Meta, error) {
	reader, err := pdf.NewReader(d.r, d.size)
	if err != nil {
		return nil, nil, err
	}

	page := reader.Page(pdfPageNum)
	chunks := page.Content().Text

	// Decode meta.
	meta, size := d.decodeMeta(chunks)
	chunks = chunks[size:]

	// Decode units.
	units := make([]schedule.Unit, 0)
	for len(chunks) > 0 {
		unit, size := d.decodeUnit(chunks)
		units = append(units, unit)
		chunks = chunks[size:]
	}

	return units, meta, nil
}
