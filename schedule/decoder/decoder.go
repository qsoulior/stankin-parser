package decoder

import "github.com/qsoulior/stankin-parser/schedule"

type Decoder interface {
	Decode() ([]schedule.Cell, schedule.Meta)
}
