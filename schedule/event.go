package schedule

import (
	"time"
)

type EventType string

const (
	EventTypeLecture EventType = "lecture"
	EventTypeSeminar EventType = "seminar"
	EventTypeLab     EventType = "lab"
)

type EventSubgroup string

const (
	EventSubgroupA EventSubgroup = "A"
	EventSubgroupB EventSubgroup = "B"
)

type Event struct {
	Title    string        `json:"title"`
	Teacher  string        `json:"teacher"`
	Type     EventType     `json:"type"`
	Subgroup EventSubgroup `json:"subgroup"`
	Location string        `json:"location"`
	Dates    []EventDate   `json:"dates"`
}

type EventFrequency string

const (
	EventFrequencyOnce     EventFrequency = "once"
	EventFrequencyWeekly   EventFrequency = "weekly"
	EventFrequencyBiweekly EventFrequency = "biweekly"
)

type EventDate struct {
	Start     time.Time      `json:"start"`
	End       time.Time      `json:"end"`
	Frequency EventFrequency `json:"frequency"`
}
