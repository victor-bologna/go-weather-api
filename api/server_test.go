package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/victor-bologna/go-weather-api/storage"
	"github.com/victor-bologna/go-weather-api/types"
)

var weather = types.WeatherModel{
	Latitude:  -35.25,
	Longitude: -30.375,
	Timezone:  "America/Sao_Paulo",
	CurrentUnits: types.CurrentUnits{
		Time:          "iso8601",
		Interval:      "seconds",
		Temperature2M: "°C",
	},
	Current: types.Current{
		Time:          "2024-03-10T14:30",
		Interval:      900,
		Temperature2M: 27.5,
	},
}

func TestGetWeather_Success(t *testing.T) {
	storageMock := &storage.MockStorage{}
	matchedBy := mock.MatchedBy(func(w *types.WeatherModel) bool {
		if w.Latitude == weather.Latitude &&
			w.Longitude == weather.Longitude {
			return true
		}
		return false
	})
	storageMock.On("SaveWeather", matchedBy).Return(nil)

	server := NewWeatherAPI(":8080", storageMock, http.NewServeMux())
	server.InitializeHandlers()

	resp := executeRequest(server, "?latitude=-35.25&longitude=-30.375")

	expectedBody := `{"latitude":-35.25,"longitude":-30.375,"timezone":"America/Sao_Paulo"`

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), expectedBody)
}

func TestGetWeather_MissingParams(t *testing.T) {
	server := NewWeatherAPI(":8080", nil, http.NewServeMux())
	server.InitializeHandlers()

	resp := executeRequest(server, "")

	errorMsg := "Latitude param not found.Longitude param not found."

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), errorMsg)
}

func executeRequest(server *WeatherServer, params string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weather"+params, nil)
	server.mux.ServeHTTP(rr, req)
	return rr
}
