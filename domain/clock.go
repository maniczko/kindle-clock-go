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

func NewSystemClock() *SystemClock {
	timezoneName := os.Getenv("APP_TIMEZONE")
	if timezoneName == "" {
		timezoneName = defaultTimezone
	}

	location, err := time.LoadLocation(timezoneName)
	if err != nil {
		panic("load APP_TIMEZONE: " + err.Error())
	}

	return &SystemClock{timezone: location}
}

func (c *SystemClock) Now() time.Time {
	return time.Now().In(c.timezone)
}
