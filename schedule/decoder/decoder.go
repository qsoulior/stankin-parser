// Package decoder provides interface and methods to decode schedule units.
package decoder

import "github.com/qsoulior/stankin-parser/schedule"

// Decoder represents abstract decoder that can decode schedule units.
type Decoder interface {
	// Decode decodes schedule units and metadata and returns them.
	Decode() ([]schedule.Unit, *schedule.Meta, error)
}
