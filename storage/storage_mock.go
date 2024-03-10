package storage

import (
	"github.com/stretchr/testify/mock"
	"github.com/victor-bologna/go-weather-api/types"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SaveWeather(weather *types.WeatherModel) error {
	args := m.Called(weather)
	return args.Error(0)
}
