package ical_encoder

import "time"

const (
	version           = "2.0"
	method            = "PUBLISH"
	productID         = "-//Unknown//Stankin parser//RU"
	timezoneID        = "Europe/Moscow"
	timezoneOffset    = "+0300"
	scale             = "GREGORIAN"
	transparent       = "OPAQUE"
	layoutDatetime    = "20060102T150405"
	layoutDatetimeUTC = "20060102T150405Z"
)

var dateWeekdays = map[time.Weekday]string{
	time.Sunday:    "SU",
	time.Monday:    "MO",
	time.Tuesday:   "TU",
	time.Wednesday: "WE",
	time.Thursday:  "TH",
	time.Friday:    "FR",
	time.Saturday:  "SA",
}
