package encoder

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/qsoulior/stankin-parser/schedule"
)

const (
	icalVersion           = "2.0"
	icalMethod            = "PUBLISH"
	icalProductID         = "-//Unknown//Stankin parser//RU"
	icalTimezoneID        = "Europe/Moscow"
	icalTimezoneOffset    = "+0300"
	icalScale             = "GREGORIAN"
	icalTransparent       = "OPAQUE"
	icalLayoutDatetime    = "20060102T150405"
	icalLayoutDatetimeUTC = "20060102T150405Z"
)

type icalEncoder struct {
	w io.Writer
}

// NewIcal creates and returns new ical encoder.
func NewIcal(w io.Writer) Encoder { return &icalEncoder{w} }

// Encode encodes schedule events and group into ical format.
// It uses subgroup to filter events.
func (e *icalEncoder) Encode(events []schedule.Event, group string, subgroup schedule.EventSubgroup) error {
	fmt.Fprint(e.w, "BEGIN:VCALENDAR\r\n")

	// VCALENDAR
	fmt.Fprintf(e.w, "PRODID:%s\r\n", icalProductID)
	fmt.Fprintf(e.w, "VERSION:%s\r\n", icalVersion)
	fmt.Fprintf(e.w, "METHOD:%s\r\n", icalMethod)
	fmt.Fprintf(e.w, "CALSCALE:%s\r\n", icalScale)

	// VTIMEZONE
	fmt.Fprint(e.w, "BEGIN:VTIMEZONE\r\n")
	fmt.Fprintf(e.w, "TZID:%s\r\n", icalTimezoneID)
	fmt.Fprint(e.w, "BEGIN:STANDARD\r\n")
	fmt.Fprint(e.w, "DTSTART:19700101T000000\r\n")
	fmt.Fprintf(e.w, "TZOFFSETFROM:%s\r\n", icalTimezoneOffset)
	fmt.Fprintf(e.w, "TZOFFSETTO:%s\r\n", icalTimezoneOffset)
	fmt.Fprint(e.w, "END:STANDARD\r\n")
	fmt.Fprint(e.w, "END:VTIMEZONE\r\n")

	// VEVENT
	for _, event := range events {
		if subgroup == "" || event.Subgroup == "" || subgroup == event.Subgroup {
			e.encodeEvent(event)
			fmt.Fprint(e.w, "\r\n")
		}
	}

	fmt.Fprint(e.w, "END:VCALENDAR\r\n")
	return nil
}

func (e *icalEncoder) encodeEvent(event schedule.Event) {
	fmt.Fprint(e.w, "BEGIN:VEVENT\r\n")

	fmt.Fprintf(e.w, "UID:%s\r\n", uuid.New())
	fmt.Fprintf(e.w, "DTSTAMP:%s\r\n", time.Now().UTC().Format(icalLayoutDatetimeUTC))
	fmt.Fprintf(e.w, "LOCATION:%s\r\n", event.Location)
	fmt.Fprintf(e.w, "TRANSP:%s\r\n", icalTransparent)

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
		icalTimezoneID, event.Dates[0].Start.Format(icalLayoutDatetime))
	fmt.Fprintf(e.w, "DTEND;TZID=%s:%s\r\n",
		icalTimezoneID, event.Dates[0].End.Format(icalLayoutDatetime))

	dates := make([]string, len(event.Dates))
	for i, date := range event.Dates {
		dates[i] = date.Start.Format(icalLayoutDatetime)
	}

	fmt.Fprintf(e.w, "RDATE;TZID=%s:%s\r\n",
		icalTimezoneID, strings.Join(dates, ","))

	fmt.Fprint(e.w, "END:VEVENT")
}
