package config

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sethvargo/go-envconfig"
)

// DisplayConfiguration describes the final PNG returned to a Kindle.
// Rotation 90 is the default for a Kindle mounted in landscape orientation.
type DisplayConfiguration struct {
	Width    int `env:"DISPLAY_WIDTH,default=600"`
	Height   int `env:"DISPLAY_HEIGHT,default=800"`
	Rotation int `env:"DISPLAY_ROTATION,default=90"`
}

func NewDisplayConfiguration(ctx context.Context) *DisplayConfiguration {
	var c DisplayConfiguration
	if err := envconfig.Process(ctx, &c); err != nil {
		slog.Error("failed to process display configuration", "error", err)
		panic(err)
	}
	if err := c.Validate(); err != nil {
		panic(err)
	}
	return &c
}

func (c DisplayConfiguration) Validate() error {
	if c.Width <= 0 || c.Height <= 0 {
		return fmt.Errorf("DISPLAY_WIDTH and DISPLAY_HEIGHT must be positive")
	}
	if c.Rotation != 0 && c.Rotation != 90 {
		return fmt.Errorf("DISPLAY_ROTATION must be 0 or 90")
	}
	return nil
}
