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

func parseTimeStart(pos int) (time.Duration, error) {
	switch {
	case pos >= 701:
		return 21*time.Hour + 20*time.Minute, nil
	case pos >= 607:
		return 19*time.Hour + 40*time.Minute, nil
	case pos >= 514:
		return 18 * time.Hour, nil
	case pos >= 420:
		return 16 * time.Hour, nil
	case pos >= 327:
		return 14*time.Hour + 10*time.Minute, nil
	case pos >= 233:
		return 12*time.Hour + 20*time.Minute, nil
	case pos >= 139:
		return 10*time.Hour + 20*time.Minute, nil
	case pos >= 46:
		return 8*time.Hour + 30*time.Minute, nil
	default:
		return 0, errors.New("event start is invalid")
	}
}

func parseTimeEnd(pos int) (time.Duration, error) {
	switch {
	case pos < 46:
		return 0, errors.New("event end is invalid")
	case pos < 139:
		return 10*time.Hour + 10*time.Minute, nil
	case pos < 233:
		return 12 * time.Hour, nil
	case pos < 327:
		return 14 * time.Hour, nil
	case pos < 420:
		return 15*time.Hour + 50*time.Minute, nil
	case pos < 514:
		return 17*time.Hour + 40*time.Minute, nil
	case pos < 607:
		return 19*time.Hour + 30*time.Minute, nil
	case pos < 701:
		return 21*time.Hour + 10*time.Minute, nil
	default:
		return 22*time.Hour + 50*time.Minute, nil
	}
}

func parseTime(left, right int) (EventTime, error) {
	var time EventTime
	var err error

	time.Start, err = parseTimeStart(left)
	if err != nil {
		return time, err
	}

	time.End, err = parseTimeEnd(right)
	if err != nil {
		return time, err
	}

	return time, nil
}

const dateLayout = "02.01"

var dateLocation = time.FixedZone("UTC+3", 3*60*60)

var dateWeeks = map[string]int{
	"":     0,
	"к.н.": 1,
	"ч.н.": 2,
}

func parseDate(data string, year int, t EventTime) ([]EventDate, error) {
	parts := make([]string, 3)
	copy(parts, strings.FieldsFunc(data, func(r rune) bool {
		return r == '-' || r == ' '
	}))

	weeks, ok := dateWeeks[parts[2]]
	if !ok {
		return nil, ErrEventFrequencyInvalid
	}

	start, err := time.ParseInLocation(dateLayout, parts[0], dateLocation)
	if err != nil {
		return nil, ErrEventDateInvalid
	}
	start = start.AddDate(year, 0, 0)

	if weeks > 0 {
		end, err := time.ParseInLocation(dateLayout, parts[1], dateLocation)
		if err != nil {
			return nil, ErrEventDateInvalid
		}
		end = end.AddDate(year, 0, 0)

		dates := make([]EventDate, 0)
		for date := start; !date.After(end); date = date.AddDate(0, 0, 7*weeks) {
			dates = append(dates, EventDate{
				Start: date.Add(t.Start),
				End:   date.Add(t.End)},
			)
		}

		return dates, nil
	}

	date := EventDate{
		Start: start.Add(t.Start),
		End:   start.Add(t.End),
	}

	return []EventDate{date}, nil
}

func parse(unit Unit, year int) (Event, error) {
	var event Event

	// Parse type.
	re := regexp.MustCompile(
		fmt.Sprintf(` (%s|%s|%s)\. `, EventTypeLecture, EventTypeSeminar, EventTypeLab),
	)

	loc := re.FindStringSubmatchIndex(unit.Data)
	if len(loc) < 4 {
		return event, ErrEventTypeNotFound
	}

	event.Type = EventType(unit.Data[loc[2]:loc[3]])

	// "<title>" OR "<title>", "<teacher>"
	prefix := strings.Split(unit.Data[:loc[0]], ". ")
	if n := len(prefix); n < 1 || n > 2 {
		return event, ErrEventInvalid
	}

	// "<location>", "[<date>]" OR "(<subgroup>)", "<location>", "[<date>, <date>]"
	suffix := strings.Split(unit.Data[loc[1]:], ". ")
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
	time, err := parseTime(unit.Left, unit.Right)
	if err != nil {
		return event, err
	}

	// Parse dates.
	parts := strings.Split(strings.Trim(suffix[len(suffix)-1], "[]"), ", ")
	for _, part := range parts {
		dates, err := parseDate(part, year, time)
		if err != nil {
			return event, err
		}
		event.Dates = append(event.Dates, dates...)
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

// Parse schedule units and returns events.
// Since any unit does not contain year of event, it is accepted as required argument.
func Parse(units []Unit, year int) ([]Event, error) {
	var err error

	events := make([]Event, len(units))
	for i, unit := range units {
		events[i], err = parse(unit, year)
		if err != nil {
			return nil, err
		}
	}

	return events, nil
}
