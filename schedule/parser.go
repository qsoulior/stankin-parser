package schedule

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	ErrEventInvalid          = errors.New("event is invalid")
	ErrEventTypeNotFound     = errors.New("event type is not found")
	ErrEventDateInvalid      = errors.New("event date is invalid")
	ErrEventFrequencyInvalid = errors.New("event frequency is invalid")
)

type EventClock struct {
	hour int
	min  int
}

type EventTime struct {
	start EventClock
	end   EventClock
}

func parseTimeStart(pos int) (EventClock, error) {
	switch {
	case pos >= 701:
		return EventClock{21, 20}, nil
	case pos >= 607:
		return EventClock{19, 40}, nil
	case pos >= 514:
		return EventClock{18, 0}, nil
	case pos >= 420:
		return EventClock{16, 0}, nil
	case pos >= 327:
		return EventClock{14, 10}, nil
	case pos >= 233:
		return EventClock{12, 20}, nil
	case pos >= 139:
		return EventClock{10, 20}, nil
	case pos >= 46:
		return EventClock{8, 30}, nil
	default:
		return EventClock{}, errors.New("event start is invalid")
	}
}

func parseTimeEnd(pos int) (EventClock, error) {
	switch {
	case pos < 46:
		return EventClock{}, errors.New("event end is invalid")
	case pos < 139:
		return EventClock{10, 10}, nil
	case pos < 233:
		return EventClock{12, 0}, nil
	case pos < 327:
		return EventClock{14, 0}, nil
	case pos < 420:
		return EventClock{15, 50}, nil
	case pos < 514:
		return EventClock{17, 40}, nil
	case pos < 607:
		return EventClock{19, 30}, nil
	case pos < 701:
		return EventClock{21, 10}, nil
	default:
		return EventClock{22, 50}, nil
	}
}

func parseTime(left, right int) (EventTime, error) {
	var time EventTime
	var err error

	time.start, err = parseTimeStart(left)
	if err != nil {
		return time, err
	}

	time.end, err = parseTimeEnd(right)
	if err != nil {
		return time, err
	}

	return time, nil
}

const dateLayout = "02.01"

var dateLocation = time.FixedZone("UTC+3", 3*60*60)

func parseDate(data string, year int, offset EventTime) (EventDate, error) {
	var (
		date       EventDate
		err        error
		start, end string
	)

	splitDate := strings.FieldsFunc(data, func(r rune) bool {
		return r == '-' || r == ' '
	})

	switch len(splitDate) {
	case 1:
		date.Frequency = EventFrequencyOnce
		start = splitDate[0]
		end = splitDate[0]
	case 3:
		switch splitDate[2] {
		case "к.н.":
			date.Frequency = EventFrequencyWeekly
		case "ч.н.":
			date.Frequency = EventFrequencyBiweekly
		default:
			return date, ErrEventFrequencyInvalid
		}

		start = splitDate[0]
		end = splitDate[1]
	default:
		return date, ErrEventDateInvalid
	}

	date.Start, err = time.ParseInLocation(dateLayout, start, dateLocation)
	if err != nil {
		return date, ErrEventDateInvalid
	}

	date.End, err = time.ParseInLocation(dateLayout, end, dateLocation)
	if err != nil {
		return date, ErrEventDateInvalid
	}

	// Normalize date.
	date.Start = date.Start.
		AddDate(year, 0, 0).
		Add(time.Hour*time.Duration(offset.start.hour) + time.Minute*time.Duration(offset.start.min))

	date.End = date.End.
		AddDate(year, 0, 0).
		Add(time.Hour*time.Duration(offset.end.hour) + time.Minute*time.Duration(offset.end.min))

	return date, nil
}

func parse(cell Cell, year int) (Event, error) {
	var event Event

	// Parse type.
	re := regexp.MustCompile(
		fmt.Sprintf(` (%s|%s|%s)\. `, EventTypeLecture, EventTypeSeminar, EventTypeLab),
	)

	loc := re.FindStringSubmatchIndex(cell.Data)
	if len(loc) < 4 {
		return event, ErrEventTypeNotFound
	}

	event.Type = EventType(cell.Data[loc[2]:loc[3]])

	// "<title>" OR "<title>", "<teacher>"
	prefix := strings.Split(cell.Data[:loc[0]], ". ")
	if n := len(prefix); n < 1 || n > 2 {
		return event, ErrEventInvalid
	}

	// "<location>", "[<date>]" OR "(<subgroup>)", "<location>", "[<date>, <date>]"
	suffix := strings.Split(cell.Data[loc[1]:], ". ")
	if n := len(suffix); n < 2 || n > 3 {
		return event, ErrEventInvalid
	}

	// Parse title and teacher.
	if len(prefix) == 1 {
		event.Title = prefix[0][:len(prefix[0])-1]
	} else {
		event.Title = prefix[0]
		event.Teacher = prefix[1]
	}

	// Parse time.
	time, err := parseTime(cell.Left, cell.Right)
	if err != nil {
		return event, err
	}

	// Parse dates.
	parts := strings.Split(strings.Trim(suffix[len(suffix)-1], "[]"), ", ")
	event.Dates = make([]EventDate, len(parts))

	for i, part := range parts {
		event.Dates[i], err = parseDate(part, year, time)
		if err != nil {
			return event, err
		}
	}

	// Parse subgroup and location.
	if len(suffix) == 2 {
		event.Location = suffix[0]
	} else {
		event.Subgroup = EventSubgroup(strings.Trim(suffix[0], "()"))
		event.Location = suffix[1]
	}

	return event, nil
}

func Parse(cells []Cell, year int) ([]Event, error) {
	var err error

	events := make([]Event, len(cells))
	for i, cell := range cells {
		events[i], err = parse(cell, year)
		if err != nil {
			return nil, err
		}
	}

	return events, nil
}
