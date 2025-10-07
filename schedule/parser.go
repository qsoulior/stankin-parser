package schedule

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var location = time.FixedZone("UTC+3", 3*60*60)

var (
	ErrEventInvalid         = errors.New("event is invalid")
	ErrEventTypeNotFound    = errors.New("event type is not found")
	ErrEventDateInvalid     = errors.New("event date is invalid")
	ErrEventIntervalInvalid = errors.New("event interval is invalid")
	ErrEventStartInvalid    = errors.New("event start is invalid")
	ErrEventEndInvalid      = errors.New("event end is invalid")
)

func parseTimeStart(pos int) (time.Time, error) {
	switch {
	case pos >= 701:
		return time.Date(0, 0, 0, 21, 20, 0, 0, location), nil
	case pos >= 607:
		return time.Date(0, 0, 0, 19, 40, 0, 0, location), nil
	case pos >= 514:
		return time.Date(0, 0, 0, 18, 00, 0, 0, location), nil
	case pos >= 420:
		return time.Date(0, 0, 0, 16, 00, 0, 0, location), nil
	case pos >= 327:
		return time.Date(0, 0, 0, 14, 10, 0, 0, location), nil
	case pos >= 233:
		return time.Date(0, 0, 0, 12, 20, 0, 0, location), nil
	case pos >= 139:
		return time.Date(0, 0, 0, 10, 20, 0, 0, location), nil
	case pos >= 46:
		return time.Date(0, 0, 0, 8, 30, 0, 0, location), nil
	default:
		return time.Time{}, ErrEventStartInvalid
	}
}

func parseTimeEnd(pos int) (time.Time, error) {
	switch {
	case pos < 46:
		return time.Time{}, ErrEventEndInvalid
	case pos < 139:
		return time.Date(0, 0, 0, 10, 10, 0, 0, location), nil
	case pos < 233:
		return time.Date(0, 0, 0, 12, 0, 0, 0, location), nil
	case pos < 327:
		return time.Date(0, 0, 0, 14, 0, 0, 0, location), nil
	case pos < 420:
		return time.Date(0, 0, 0, 15, 50, 0, 0, location), nil
	case pos < 514:
		return time.Date(0, 0, 0, 17, 40, 0, 0, location), nil
	case pos < 607:
		return time.Date(0, 0, 0, 19, 30, 0, 0, location), nil
	case pos < 701:
		return time.Date(0, 0, 0, 21, 10, 0, 0, location), nil
	default:
		return time.Date(0, 0, 0, 22, 50, 0, 0, location), nil
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

var dateIntervals = map[string]int{
	"":     0,
	"к.н.": 1,
	"ч.н.": 2,
}

func parseDate(data string, year int) (EventDate, error) {
	parts := make([]string, 3)
	copy(parts, strings.FieldsFunc(data, func(r rune) bool {
		return r == '-' || r == ' '
	}))

	interval, ok := dateIntervals[parts[2]]
	if !ok {
		return EventDate{}, ErrEventIntervalInvalid
	}

	start, err := time.ParseInLocation(dateLayout, parts[0], location)
	if err != nil {
		return EventDate{}, ErrEventDateInvalid
	}
	start = start.AddDate(year, 0, 0)

	if interval == 0 {
		date := EventDate{
			Start:    start,
			End:      start,
			Interval: 0,
		}

		return date, nil
	}

	end, err := time.ParseInLocation(dateLayout, parts[1], location)
	if err != nil {
		return EventDate{}, ErrEventDateInvalid
	}
	end = end.AddDate(year, 0, 0)

	date := EventDate{
		Start:    start,
		End:      end,
		Interval: interval,
	}

	return date, nil
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
	event.Time = time

	// Parse dates.
	parts := strings.Split(strings.Trim(suffix[len(suffix)-1], "[]"), ", ")
	for _, part := range parts {
		date, err := parseDate(part, year)
		if err != nil {
			return event, err
		}
		event.Dates = append(event.Dates, date)
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
