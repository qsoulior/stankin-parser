package encoder

import "github.com/qsoulior/stankin-parser/schedule"

type Encoder interface {
	Encode(events []schedule.Event, group string, subgroup schedule.EventSubgroup)
}
