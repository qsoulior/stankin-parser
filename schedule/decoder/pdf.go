package decoder

import (
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/qsoulior/stankin-parser/schedule"
)

const PageNum = 1

type Chunk struct {
	Data string
	X    int
	Y    int
}

func decodeMeta(chunks []Chunk) (*schedule.Meta, int) {
	i := 1
	for chunks[i].Y == chunks[i-1].Y && i < len(chunks) {
		i++
	}

	var data strings.Builder
	data.Grow(i - 1)
	for j := 0; j < i; j++ {
		data.WriteString(chunks[j].Data)
	}

	meta := schedule.Meta{
		Group: data.String(),
	}

	for chunks[i].X < 42 || chunks[i].Y > 520 && i < len(chunks) {
		i++
	}

	return &meta, i
}

func decodeCell(chunks []Chunk) (schedule.Cell, int) {
	var data strings.Builder
	data.Grow(len(chunks[0].Data))
	data.WriteString(chunks[0].Data)

	i := 1
	iMax := 0
	for chunks[i-1].Data != "]" && i < len(chunks) {
		if chunks[i].X > chunks[iMax].X {
			iMax = i
		}

		if chunks[i].Y != chunks[i-1].Y {
			data.WriteRune(' ')
		}
		data.WriteString(chunks[i].Data)
		i++
	}

	cell := schedule.Cell{
		Data:   data.String(),
		Left:   chunks[0].X,
		Top:    chunks[0].Y,
		Right:  chunks[iMax].X,
		Bottom: chunks[i-1].Y,
	}

	return cell, i
}

type pdfe struct {
	r    io.ReaderAt
	size int64
}

func NewPDF(r io.ReaderAt, size int64) *pdfe {
	return &pdfe{r, size}
}

func (d *pdfe) Decode() ([]schedule.Cell, *schedule.Meta, error) {
	reader, err := pdf.NewReader(d.r, d.size)
	if err != nil {
		return nil, nil, err
	}

	page := reader.Page(PageNum)
	texts := page.Content().Text
	chunks := make([]Chunk, len(texts))

	for i, text := range texts {
		chunks[i] = Chunk{
			Data: text.S,
			X:    int(text.X),
			Y:    int(text.Y),
		}
	}

	// Decode meta.
	meta, size := decodeMeta(chunks)
	chunks = chunks[size:]

	// Decode cells.
	cells := make([]schedule.Cell, 0)
	for len(chunks) > 0 {
		cell, size := decodeCell(chunks)
		cells = append(cells, cell)
		chunks = chunks[size:]
	}

	return cells, meta, nil
}
