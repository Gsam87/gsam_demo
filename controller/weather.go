package controller

import (
	repo_entity "github.com/qww83728/gsam_demo/domain/entity/repo"
	repo "github.com/qww83728/gsam_demo/domain/repository"
)

type WeatherController interface {
	GetTodayInfo() ([]repo_entity.Weather, error)
}

type WeatherControllerImpl struct {
	weatherRepo repo.WeatherRepo
}

func NewWeatherController(
	weatherRepo repo.WeatherRepo,
) WeatherController {
	return &WeatherControllerImpl{
		weatherRepo: weatherRepo,
	}
}

func (c *WeatherControllerImpl) GetTodayInfo() ([]repo_entity.Weather, error) {
	weathers, err := c.weatherRepo.GetTodayInfo()
	if err != nil {
		return nil, err
	}

	return weathers, nil
}
