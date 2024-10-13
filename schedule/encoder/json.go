package encoder

import (
	"encoding/json"
	"io"

	"github.com/qsoulior/stankin-parser/schedule"
)

type jsone struct {
	w io.Writer
}

func NewJSON(w io.Writer) *jsone { return &jsone{w} }

func (e *jsone) Encode(events []schedule.Event, group string, subgroup schedule.EventSubgroup) error {
	enc := json.NewEncoder(e.w)

	return enc.Encode(map[string]any{
		"group":    group,
		"subgroup": subgroup,
		"events":   events,
	})
}
