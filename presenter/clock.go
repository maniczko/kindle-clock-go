package presenter

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/y-yu/kindle-clock-go/config"
	"github.com/y-yu/kindle-clock-go/domain"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type ClockHandler struct {
	display config.DisplayConfiguration
	font    *truetype.Font
	clock   domain.Clock
}

func NewClockHandler(fontConfig *config.FontConfiguration, displayConfig *config.DisplayConfiguration, clock domain.Clock) *ClockHandler {
	fontFile, err := os.ReadFile(fontConfig.DosisFontPath)
	if err != nil {
		panic(fmt.Sprintf("load clock font %q: %v", fontConfig.DosisFontPath, err))
	}

	parsedFont, err := truetype.Parse(fontFile)
	if err != nil {
		panic(fmt.Sprintf("parse clock font %q: %v", fontConfig.DosisFontPath, err))
	}

	return NewClockHandlerFromFont(parsedFont, *displayConfig, clock)
}

// NewClockHandlerFromFont keeps rendering testable without a font file on disk.
func NewClockHandlerFromFont(font *truetype.Font, display config.DisplayConfiguration, clock domain.Clock) *ClockHandler {
	return &ClockHandler{display: display, font: font, clock: clock}
}

func (h *ClockHandler) Handle(w http.ResponseWriter, r *http.Request) {
	buf, err := h.generatePNG()
	if err != nil {
		slog.Error("failed to create clock PNG", "error", err)
		http.Error(w, "failed to generate clock image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(buf.Bytes()); err != nil {
		slog.Error("failed to write clock PNG", "error", err)
	}
}

func (h *ClockHandler) generatePNG() (bytes.Buffer, error) {
	if err := h.display.Validate(); err != nil {
		return bytes.Buffer{}, err
	}

	canvasWidth, canvasHeight := h.display.Width, h.display.Height
	if h.display.Rotation == 90 {
		canvasWidth, canvasHeight = canvasHeight, canvasWidth
	}

	img := image.NewGray(image.Rect(0, 0, canvasWidth, canvasHeight))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	now := h.clock.Now()

	h.drawCentered(img, now.Format("02.01.2006"), float64(canvasHeight)*0.18, float64(canvasHeight)*0.23)
	h.drawCentered(img, now.Format("15:04"), float64(canvasHeight)*0.46, float64(canvasHeight)*0.72)

	result := image.Image(img)
	if h.display.Rotation == 90 {
		result = rotate90(img)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, result); err != nil {
		return buf, err
	}
	return buf, nil
}

func (h *ClockHandler) drawCentered(dst draw.Image, value string, size, baseline float64) {
	face := truetype.NewFace(h.font, &truetype.Options{Size: size})
	drawer := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  fixed.P(0, int(baseline)),
	}
	DrawStringCentering(drawer, dst.Bounds().Dx(), value)
}

func rotate90(src image.Image) image.Image {
	srcBounds := src.Bounds()
	width, height := srcBounds.Dx(), srcBounds.Dy()
	dst := image.NewGray(image.Rect(0, 0, height, width))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			dst.Set(y, width-1-x, src.At(srcBounds.Min.X+x, srcBounds.Min.Y+y))
		}
	}
	return dst
}
