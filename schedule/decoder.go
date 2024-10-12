package schedule

import "strings"

type Meta struct {
	Group string
}

func decodeMeta(chunks []Chunk) (Meta, int) {
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

	for chunks[i].X < 42 || chunks[i].Y > 520 && i < len(chunks) {
		i++
	}

	return meta, i
}

type Cell struct {
	Data                     string
	Left, Right, Top, Bottom int
}

func decodeCell(chunks []Chunk) (Cell, int) {
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

	cell := Cell{
		Data:   data.String(),
		Left:   chunks[0].X,
		Top:    chunks[0].Y,
		Right:  chunks[iMax].X,
		Bottom: chunks[i-1].Y,
	}

	return cell, i
}

func Decode(chunks []Chunk) ([]Cell, Meta) {
	// Decode meta.
	meta, size := decodeMeta(chunks)
	chunks = chunks[size:]

	// Decode cells.
	cells := make([]Cell, 0)
	for len(chunks) > 0 {
		cell, size := decodeCell(chunks)
		cells = append(cells, cell)
		chunks = chunks[size:]
	}

	return cells, meta
}
