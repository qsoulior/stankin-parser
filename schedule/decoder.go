package schedule

import "strings"

type Meta struct {
	Group string
}

func DecodeMeta(chunks []Chunk) (Meta, int) {
	i := 1
	for chunks[i].Y == chunks[i-1].Y && i < len(chunks) {
		i++
	}

	var data strings.Builder
	data.Grow(i - 1)
	for j := 0; j < i; j++ {
		data.WriteString(chunks[j].Data)
	}

	meta := Meta{
		Group: data.String(),
	}

	for chunks[i].X < 46 || chunks[i].Y > 510 && i < len(chunks) {
		i++
	}

	return meta, i
}

type Cell struct {
	Data                     string
	Left, Right, Top, Bottom int
}

func DecodeCell(chunks []Chunk) (Cell, int) {
	var data strings.Builder
	data.Grow(len(chunks[0].Data))
	data.WriteString(chunks[0].Data)

	i := 1
	for chunks[i-1].Data != "]" && i < len(chunks) {
		if chunks[i].Y != chunks[i-1].Y {
			data.WriteRune(' ')
		}
		data.WriteString(chunks[i].Data)
		i++
	}

	cell := Cell{
		Data:   data.String(),
		Left:   chunks[0].X,
		Top:    chunks[0].Y,
		Right:  chunks[i-1].X,
		Bottom: chunks[i-1].Y,
	}

	return cell, i
}

func DecodeData(chunks []Chunk) ([]Cell, Meta) {
	// Decode meta.
	meta, size := DecodeMeta(chunks)
	chunks = chunks[size:]

	// Decode cells.
	cells := make([]Cell, 0)
	for len(chunks) > 0 {
		cell, size := DecodeCell(chunks)
		cells = append(cells, cell)
		chunks = chunks[size:]
	}

	return cells, meta
}
