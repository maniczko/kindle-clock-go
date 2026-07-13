package presenter

import (
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/require"
	"github.com/y-yu/kindle-clock-go/config"
	"golang.org/x/image/font/gofont/goregular"
)

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

func newTestClockHandler(t *testing.T, rotation int) *ClockHandler {
	t.Helper()
	font, err := truetype.Parse(goregular.TTF)
	require.NoError(t, err)
	return NewClockHandlerFromFont(font, config.DisplayConfiguration{
		Width: 600, Height: 800, Rotation: rotation,
	}, fixedClock{now: time.Date(2026, time.July, 13, 9, 15, 0, 0, time.FixedZone("Europe/Warsaw", 2*60*60))})
}

func TestClockHandlerReturnsKindleSizedPNG(t *testing.T) {
	handler := newTestClockHandler(t, 90)
	response := httptest.NewRecorder()

	handler.Handle(response, httptest.NewRequest(http.MethodGet, "/clock", nil))

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "image/png", response.Header().Get("Content-Type"))
	require.Equal(t, "no-store, max-age=0", response.Header().Get("Cache-Control"))

	image, err := png.Decode(response.Result().Body)
	require.NoError(t, err)
	require.Equal(t, 600, image.Bounds().Dx())
	require.Equal(t, 800, image.Bounds().Dy())
	require.True(t, hasLightAndDarkPixels(image))
}

func TestClockHandlerKeepsConfiguredDimensionsWithoutRotation(t *testing.T) {
	handler := newTestClockHandler(t, 0)
	buffer, err := handler.generatePNG()
	require.NoError(t, err)

	image, err := png.Decode(&buffer)
	require.NoError(t, err)
	require.Equal(t, 600, image.Bounds().Dx())
	require.Equal(t, 800, image.Bounds().Dy())
	require.True(t, hasLightAndDarkPixels(image))
}

func TestDisplayConfigurationRejectsUnsupportedRotation(t *testing.T) {
	require.Error(t, (config.DisplayConfiguration{Width: 600, Height: 800, Rotation: 180}).Validate())
}

func hasLightAndDarkPixels(img image.Image) bool {
	var hasLight, hasDark bool
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			red, _, _, _ := img.At(x, y).RGBA()
			if red == 0 {
				hasDark = true
			}
			if red == 0xffff {
				hasLight = true
			}
		}
	}
	return hasLight && hasDark
}
