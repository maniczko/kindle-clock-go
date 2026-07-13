package domain

import (
	"os"
	"time"

	_ "time/tzdata"
)

type Clock interface {
	Now() time.Time
}

type SystemClock struct {
	timezone *time.Location
}

const defaultTimezone = "Europe/Warsaw"

// JST is retained for the legacy weather integration, which parses timestamps
// from a Tokyo-based upstream API. The standalone clock never reads it.
var JST = time.FixedZone("Asia/Tokyo", 9*60*60)

func NewSystemClock() *SystemClock {
	return &SystemClock{timezone: Location()}
}

// Location returns the configured IANA timezone. It is shared by the optional
// weather integration so timestamps do not silently revert to a fixed offset.
func Location() *time.Location {
	timezoneName := os.Getenv("APP_TIMEZONE")
	if timezoneName == "" {
		timezoneName = defaultTimezone
	}

	location, err := time.LoadLocation(timezoneName)
	if err != nil {
		panic("load APP_TIMEZONE: " + err.Error())
	}

	return location
}

func (c *SystemClock) Now() time.Time {
	return time.Now().In(c.timezone)
}
