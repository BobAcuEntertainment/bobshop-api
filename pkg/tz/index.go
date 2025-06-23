package tz

import "time"

type Timezone string

const (
	UTC           Timezone = "UTC"
	Local         Timezone = "Local"
	AsiaHoChiMinh Timezone = "Asia/Ho_Chi_Minh"
)

func Location(name Timezone) *time.Location {
	rs, _ := time.LoadLocation(string(name))
	return rs
}

func Parse(layout string, value string, timezone Timezone) time.Time {
	rs, _ := time.ParseInLocation(layout, value, Location(timezone))
	return rs
}

func ParseDate(value string, timezone Timezone) time.Time {
	rs, _ := time.ParseInLocation(time.DateOnly, value, Location(timezone))
	return rs
}

func ParseDateTime(value string, timezone Timezone) time.Time {
	rs, _ := time.ParseInLocation(time.DateTime, value, Location(timezone))
	return rs
}

func StartOfDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

func EndOfDay(date time.Time) time.Time {
	return StartOfDay(date).AddDate(0, 0, 1).Add(-time.Second)
}

func StartOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

func EndOfMonth(date time.Time) time.Time {
	return StartOfMonth(date).AddDate(0, 1, 0).Add(-time.Second)
}
