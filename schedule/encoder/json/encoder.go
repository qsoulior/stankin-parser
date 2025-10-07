package json_encoder

import (
	"encoding/json"
	"io"

	"github.com/qsoulior/stankin-parser/schedule"
)

type Encoder struct {
	w io.Writer
}

// New creates and returns new json encoder.
func New(w io.Writer) *Encoder { return &Encoder{w} }

// Encode encodes schedule events and group into json format.
// It uses subgroup to filter events.
func (e *Encoder) Encode(events []schedule.Event, group string, subgroup schedule.EventSubgroup) error {
	encoder := json.NewEncoder(e.w)

	results := make([]Event, 0)
	for _, event := range events {
		if subgroup == "" || event.Subgroup == "" || subgroup == event.Subgroup {
			results = append(results, e.encodeEvent(event))
		}
	}

	return encoder.Encode(map[string]any{
		"group":    group,
		"subgroup": subgroup,
		"events":   results,
	})
}

func (e *Encoder) encodeEvent(event schedule.Event) Event {
	dates := make([]EventDate, len(event.Dates))
	for i := range event.Dates {
		dates[i] = e.encodeEventDate(event.Dates[i])
	}

	return Event{
		Title:    event.Title,
		Teacher:  event.Teacher,
		Type:     string(event.Type),
		Subgroup: string(event.Subgroup),
		Location: event.Location,
		Time:     e.encodeEventTime(event.Time),
		Dates:    dates,
	}
}

func (e *Encoder) encodeEventDate(ed schedule.EventDate) EventDate {
	return EventDate{
		Start:    ed.Start.Format(layoutDate),
		End:      ed.End.Format(layoutDate),
		Interval: ed.Interval,
	}
}

func (e *Encoder) encodeEventTime(et schedule.EventTime) EventTime {
	return EventTime{
		Start: et.Start.Format(layoutTime),
		End:   et.End.Format(layoutTime),
	}
}
