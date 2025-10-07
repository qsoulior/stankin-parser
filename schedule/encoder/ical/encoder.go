package ical_encoder

import (
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/qsoulior/stankin-parser/schedule"
)

type Encoder struct {
	w io.Writer
}

// New creates and returns new iCal encoder.
func New(w io.Writer) *Encoder { return &Encoder{w} }

// Encode encodes schedule events and group into iCal format.
// It uses subgroup to filter events.
func (e *Encoder) Encode(events []schedule.Event, group string, subgroup schedule.EventSubgroup) error {
	fmt.Fprint(e.w, "BEGIN:VCALENDAR\r\n")

	// VCALENDAR
	fmt.Fprintf(e.w, "PRODID:%s\r\n", productID)
	fmt.Fprintf(e.w, "VERSION:%s\r\n", version)
	fmt.Fprintf(e.w, "METHOD:%s\r\n", method)
	fmt.Fprintf(e.w, "CALSCALE:%s\r\n", scale)

	// VTIMEZONE
	fmt.Fprint(e.w, "BEGIN:VTIMEZONE\r\n")
	fmt.Fprintf(e.w, "TZID:%s\r\n", timezoneID)
	fmt.Fprint(e.w, "BEGIN:STANDARD\r\n")
	fmt.Fprint(e.w, "DTSTART:19700101T000000\r\n")
	fmt.Fprintf(e.w, "TZOFFSETFROM:%s\r\n", timezoneOffset)
	fmt.Fprintf(e.w, "TZOFFSETTO:%s\r\n", timezoneOffset)
	fmt.Fprint(e.w, "END:STANDARD\r\n")
	fmt.Fprint(e.w, "END:VTIMEZONE\r\n")

	// VEVENT
	for _, event := range events {
		if subgroup == "" || event.Subgroup == "" || subgroup == event.Subgroup {
			e.encodeEvent(event)
		}
	}

	fmt.Fprint(e.w, "END:VCALENDAR\r\n")
	return nil
}

func (e *Encoder) encodeEvent(event schedule.Event) {
	for _, date := range event.Dates {
		fmt.Fprint(e.w, "BEGIN:VEVENT\r\n")

		fmt.Fprintf(e.w, "UID:%s\r\n", uuid.New())
		fmt.Fprintf(e.w, "DTSTAMP:%s\r\n", time.Now().UTC().Format(layoutDatetimeUTC))
		fmt.Fprintf(e.w, "LOCATION:%s\r\n", event.Location)
		fmt.Fprintf(e.w, "TRANSP:%s\r\n", transparent)

		fmt.Fprint(e.w, "SUMMARY:")
		if event.Subgroup != "" {
			fmt.Fprintf(e.w, "[%s] ", event.Subgroup)
		}
		fmt.Fprintf(e.w, "%s\r\n", event.Title)

		fmt.Fprintf(e.w, "DESCRIPTION:%s [%s", event.Teacher, event.Type)
		if event.Subgroup != "" {
			fmt.Fprintf(e.w, " - %s", event.Subgroup)
		}
		fmt.Fprint(e.w, "]\r\n")

		fmt.Fprintf(e.w, "DTSTART;TZID=%s:%s\r\n",
			timezoneID,
			toDatetime(date.Start, event.Time.Start).Format(layoutDatetime))

		fmt.Fprintf(e.w, "DTEND;TZID=%s:%s\r\n",
			timezoneID,
			toDatetime(date.Start, event.Time.End).Format(layoutDatetime))

		if date.Interval > 0 {
			fmt.Fprintf(e.w, "RRULE:FREQ=WEEKLY;INTERVAL=%d;BYDAY=%s;UNTIL=%s\r\n",
				date.Interval,
				dateWeekdays[date.Start.Weekday()],
				date.End.Add(24*time.Hour).UTC().Format(layoutDatetimeUTC))
		}

		fmt.Fprint(e.w, "END:VEVENT\r\n")
	}
}

func toDatetime(ed time.Time, et time.Time) time.Time {
	return ed.Add(time.Duration(et.Hour())*time.Hour + time.Duration(et.Minute())*time.Minute)
}
