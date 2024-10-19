package encoder

import (
	"encoding/json"
	"io"

	"github.com/qsoulior/stankin-parser/schedule"
)

type jsonEncoder struct {
	w io.Writer
}

func NewJSON(w io.Writer) Encoder { return &jsonEncoder{w} }

func (e *jsonEncoder) Encode(events []schedule.Event, group string, subgroup schedule.EventSubgroup) error {
	enc := json.NewEncoder(e.w)

	if subgroup != "" {
		tmp := make([]schedule.Event, 0)
		for _, event := range events {
			if event.Subgroup == "" || subgroup == event.Subgroup {
				tmp = append(tmp, event)
			}
		}
		events = tmp
	}

	return enc.Encode(map[string]any{
		"group":    group,
		"subgroup": subgroup,
		"events":   events,
	})
}
