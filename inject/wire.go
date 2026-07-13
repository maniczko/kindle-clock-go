//go:build wireinject
// +build wireinject

package inject

import (
	"context"
	"github.com/google/wire"
	"github.com/y-yu/kindle-clock-go/config"
	"github.com/y-yu/kindle-clock-go/domain"
	"github.com/y-yu/kindle-clock-go/infra/api"
	"github.com/y-yu/kindle-clock-go/infra/cache"
	"github.com/y-yu/kindle-clock-go/presenter"
	"github.com/y-yu/kindle-clock-go/repository"
	"github.com/y-yu/kindle-clock-go/usecase"
)

var binding = wire.NewSet(
	domain.NewSystemClock,
	config.NewAwairConfiguration,
	config.NewAuthenticationConfiguration,
	config.NewFontConfiguration,
	config.NewDisplayConfiguration,
	config.NewNatureRemoConfiguration,
	config.NewOpenWeatherMapConfiguration,
	config.NewRedisConfiguration,
	config.NewSwitchBotConfiguration,
	api.NewAwairAPIClient,
	api.NewNatureRemoAPIClient,
	api.NewSwitchBotAPIClient,
	api.NewOpenWeatherMapAPIClient,
	cache.NewAwairCacheClient,
	cache.NewSwitchBotCacheClient,
	repository.NewAwairRepository,
	repository.NewNatureRemoRepository,
	repository.NewSwitchBotRepository,
	repository.NewOpenWeatherMapRepository,
	usecase.NewGetRoomInfoUsecase,
	presenter.NewRoomInfoHandler,
	presenter.NewClockHandler,
	presenter.NewHealthHandler,
	wire.Bind(
		new(domain.Clock),
		new(*domain.SystemClock),
	),
)

func RoomInfoHandler(ctx context.Context) *presenter.RoomInfoHandler {
	wire.Build(binding)
	return nil
}

func ClockHandler(ctx context.Context) *presenter.ClockHandler {
	wire.Build(binding)
	return nil
}

func HealthHandler(ctx context.Context) *presenter.HealthHandler {
	wire.Build(binding)
	return nil
}
