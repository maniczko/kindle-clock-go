package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSystemClockUsesConfiguredTimezone(t *testing.T) {
	t.Setenv("APP_TIMEZONE", "Europe/Warsaw")
	clock := NewSystemClock()
	require.Equal(t, "Europe/Warsaw", clock.Now().Location().String())
}

func TestSystemClockRejectsInvalidTimezone(t *testing.T) {
	t.Setenv("APP_TIMEZONE", "not/a-timezone")
	require.Panics(t, func() { NewSystemClock() })
}

func TestWarsawHasSeasonalOffsets(t *testing.T) {
	location, err := time.LoadLocation("Europe/Warsaw")
	require.NoError(t, err)
	_, winterOffset := time.Date(2026, time.January, 1, 12, 0, 0, 0, location).Zone()
	_, summerOffset := time.Date(2026, time.July, 1, 12, 0, 0, 0, location).Zone()
	require.Equal(t, 3600, winterOffset)
	require.Equal(t, 7200, summerOffset)
}
