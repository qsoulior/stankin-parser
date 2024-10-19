// Package encoder provides interface and methods to encode schedule events.
package encoder

import "github.com/qsoulior/stankin-parser/schedule"

// Encoder represents abstract encoder that can encode schedule events.
type Encoder interface {
	// Encode encodes schedule events and group.
	// It uses subgroup to filter events.
	Encode(events []schedule.Event, group string, subgroup schedule.EventSubgroup) error
}
