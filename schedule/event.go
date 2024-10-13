package schedule

import (
	"time"
)

type EventType string

const (
	EventTypeLecture EventType = "лекции"
	EventTypeSeminar EventType = "семинар"
	EventTypeLab     EventType = "лабораторные занятия"
)

type EventSubgroup string

const (
	EventSubgroupA EventSubgroup = "А"
	EventSubgroupB EventSubgroup = "Б"
)

type EventTime struct {
	Start time.Duration `json:"start"`
	End   time.Duration `json:"end"`
}

type EventDate struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type Event struct {
	Title    string        `json:"title"`
	Teacher  string        `json:"teacher"`
	Type     EventType     `json:"type"`
	Subgroup EventSubgroup `json:"subgroup"`
	Location string        `json:"location"`
	Dates    []EventDate   `json:"dates"`
}
