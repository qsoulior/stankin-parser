// Package schedule provides entities and functions to parse schedule events.
package schedule

import (
	"time"
)

// EventType represents a type of schedule event.
type EventType string

const (
	EventTypeLecture EventType = "Лекция"
	EventTypeSeminar EventType = "Семинар"
	EventTypeLab     EventType = "Лабораторная"
)

// EventSubgroup represents a subgroup related to schedule event.
type EventSubgroup string

const (
	EventSubgroupA EventSubgroup = "А"
	EventSubgroupB EventSubgroup = "Б"
)

// EventTime represents time in nanoseconds for start and end of schedule event.
type EventTime struct {
	Start time.Duration `json:"start"`
	End   time.Duration `json:"end"`
}

// EventDate represents start and end dates of schedule event.
type EventDate struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Event represents entity of schedule event.
type Event struct {
	Title    string        `json:"title"`
	Teacher  string        `json:"teacher"`
	Type     EventType     `json:"type"`
	Subgroup EventSubgroup `json:"subgroup"`
	Location string        `json:"location"`
	Dates    []EventDate   `json:"dates"`
}
