package schedule

import (
	"io"

	"github.com/ledongthuc/pdf"
)

const PageNum = 1

type Chunk struct {
	Data string
	X    int
	Y    int
}

func Read(r io.ReaderAt, size int64) ([]Chunk, error) {
	reader, err := pdf.NewReader(r, size)
	if err != nil {
		return nil, err
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

	return chunks, nil
}
