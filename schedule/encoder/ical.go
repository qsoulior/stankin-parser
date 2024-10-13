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
	icalVersion     = "2.0"
	icalMethod      = "PUBLISH"
	icalProductID   = "-//Unknown//Stankin parser//RU"
	icalTimezone    = "Europe/Moscow"
	icalScale       = "GREGORIAN"
	icalTransparent = "OPAQUE"
	icalLayoutLocal = "20060102T150405"
	icalLayoutUTC   = "20060102T150405Z"
)

type icale struct {
	w io.Writer
}

func NewIcal(w io.Writer) *icale { return &icale{w} }

func (e *icale) Encode(events []schedule.Event, group string, subgroup schedule.EventSubgroup) {
	fmt.Fprint(e.w, "BEGIN:VCALENDAR\n")

	fmt.Fprintf(e.w, "PRODID:%s\n", icalProductID)
	fmt.Fprintf(e.w, "VERSION:%s\n", icalVersion)
	fmt.Fprintf(e.w, "METHOD:%s\n", icalMethod)
	fmt.Fprintf(e.w, "CALSCALE:%s\n", icalScale)
	fmt.Fprintf(e.w, "X-WR-TIMEZONE:%s\n", icalTimezone)
	fmt.Fprintf(e.w, "X-WR-CALNAME:%s\n", group)
	fmt.Fprintf(e.w, "X-WR-CALDESC:Расписание занятий %s\n", subgroup)

	for _, event := range events {
		if subgroup == "" || event.Subgroup == "" || subgroup == event.Subgroup {
			e.serializeOne(event)
			fmt.Fprint(e.w, "\n")
		}
	}

	fmt.Fprint(e.w, "END:VCALENDAR\n")
}

func (e *icale) serializeOne(event schedule.Event) {
	fmt.Fprint(e.w, "BEGIN:VEVENT\n")

	fmt.Fprintf(e.w, "UID:%s\n", uuid.New())
	fmt.Fprintf(e.w, "DTSTAMP:%s\n", time.Now().UTC().Format(icalLayoutUTC))
	fmt.Fprintf(e.w, "SUMMARY:%s\n", event.Title)
	fmt.Fprintf(e.w, "LOCATION:%s\n", event.Location)
	fmt.Fprintf(e.w, "TRANSP:%s\n", icalTransparent)

	fmt.Fprintf(e.w, "DESCRIPTION:%s [%s", event.Teacher, event.Type)
	if event.Subgroup != "" {
		fmt.Fprintf(e.w, " - %s", event.Subgroup)
	}
	fmt.Fprint(e.w, "]\n")

	fmt.Fprintf(e.w, "DTSTART;TZID=%s:%s\n",
		icalTimezone, event.Dates[0].Start.Format(icalLayoutLocal))
	fmt.Fprintf(e.w, "DTEND;TZID=%s:%s\n",
		icalTimezone, event.Dates[0].End.Format(icalLayoutLocal))

	dates := make([]string, len(event.Dates))
	for i, date := range event.Dates {
		dates[i] = date.Start.UTC().Format(icalLayoutUTC)
	}

	fmt.Fprintf(e.w, "RDATE:%s\n", strings.Join(dates, ","))

	fmt.Fprint(e.w, "END:VEVENT")
}
