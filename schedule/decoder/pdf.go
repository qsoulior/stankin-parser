package decoder

import (
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/qsoulior/stankin-parser/schedule"
)

const PdfPageNum = 1

type pdfDecoder struct {
	r    io.ReaderAt
	size int64
}

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

func (d *pdfDecoder) decodeCell(chunks []pdf.Text) (schedule.Cell, int) {
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

	cell := schedule.Cell{
		Data:   data.String(),
		Left:   int(chunks[0].X),
		Top:    int(chunks[0].Y),
		Right:  int(chunks[iMax].X),
		Bottom: int(chunks[i-1].Y),
	}

	return cell, i
}

func (d *pdfDecoder) Decode() ([]schedule.Cell, *schedule.Meta, error) {
	reader, err := pdf.NewReader(d.r, d.size)
	if err != nil {
		return nil, nil, err
	}

	page := reader.Page(PdfPageNum)
	chunks := page.Content().Text

	// Decode meta.
	meta, size := d.decodeMeta(chunks)
	chunks = chunks[size:]

	// Decode cells.
	cells := make([]schedule.Cell, 0)
	for len(chunks) > 0 {
		cell, size := d.decodeCell(chunks)
		cells = append(cells, cell)
		chunks = chunks[size:]
	}

	return cells, meta, nil
}
