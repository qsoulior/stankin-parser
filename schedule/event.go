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

// EventTime represents time for start and end of schedule event.
type EventTime struct {
	Start time.Time
	End   time.Time
}

// EventDate represents start and end dates of schedule event.
type EventDate struct {
	Start    time.Time
	End      time.Time
	Interval int
}

// Event represents entity of schedule event.
type Event struct {
	Title    string
	Teacher  string
	Type     EventType
	Subgroup EventSubgroup
	Location string
	Time     EventTime
	Dates    []EventDate
}
