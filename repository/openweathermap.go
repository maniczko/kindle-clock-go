package repository

import (
	"context"
	"errors"
	"github.com/y-yu/kindle-clock-go/domain"
	"github.com/y-yu/kindle-clock-go/domain/api"
	"github.com/y-yu/kindle-clock-go/domain/model"
	"github.com/y-yu/kindle-clock-go/domain/repository"
	"time"
)

type OpenWeatherMapRepositoryImpl struct {
	client api.OpenWeatherMapAPIClient
}

var _ repository.OpenWeatherMapRepository = (*OpenWeatherMapRepositoryImpl)(nil)

func NewOpenWeatherMapRepository(client api.OpenWeatherMapAPIClient) repository.OpenWeatherMapRepository {
	return &OpenWeatherMapRepositoryImpl{
		client: client,
	}
}

func (o *OpenWeatherMapRepositoryImpl) GetCurrentWeather(ctx context.Context) (model.Weather, error) {
	data, err := o.client.GetLatest(ctx)
	if err != nil {
		return model.Weather{}, err
	}
	if len(data.Weather) < 1 {
		return model.Weather{}, errors.New("no current weather found")
	}

	return model.Weather{
		Icon:      data.Weather[0].Icon,
		UpdatedAt: time.Unix(data.Datetime, 0).In(domain.Location()),
	}, nil
}
