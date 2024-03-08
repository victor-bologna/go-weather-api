package storage

import "github.com/victor-bologna/go-weather-api/types"

type WeatherDB interface {
	SaveWeather(*types.WeatherModel) error
}
