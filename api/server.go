package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/victor-bologna/go-weather-api/file"
	"github.com/victor-bologna/go-weather-api/storage"
	"github.com/victor-bologna/go-weather-api/types"
)

type errorResponse struct {
	ErrorMsg       string
	ErrorCause     string
	ErrorTimestamp string
}

type WeatherServer struct {
	mux     *http.ServeMux
	port    string
	storage storage.WeatherDB
}

func NewWeatherAPI(port string, storage storage.WeatherDB, mux *http.ServeMux) *WeatherServer {
	return &WeatherServer{
		mux:     mux,
		port:    port,
		storage: storage,
	}
}

func (server *WeatherServer) StartServer() {
	server.InitializeHandlers()
	http.ListenAndServe(server.port, server.mux)
}

func (server *WeatherServer) GetWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	err := server.validateParams(r)
	if err != nil {
		server.handleError(w, http.StatusBadRequest, "Missing required params.", err)
		return
	}

	resp, err := server.makeHttpRequestWithContext(r)
	if err != nil {
		server.handleError(w, http.StatusInternalServerError, "Error when requesting weather data.", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		server.handleError(w, http.StatusInternalServerError, "Error reading weather respoonse body.", err)
		return
	}
	var weather types.WeatherModel
	err = json.Unmarshal(body, &weather)
	if err != nil {
		server.handleError(w, http.StatusBadRequest, "Error when converting to JSON", err)
		return
	}

	err = server.storage.SaveWeather(&weather)
	if err != nil {
		server.handleError(w, http.StatusBadRequest, "Error when saving Weather in database.", err)
		return
	}

	err = file.RegisterLogFile(weather, "INFO: ", types.NewStandardLogFile())
	if err != nil {
		server.handleError(w, http.StatusBadRequest, "Weather saved, but could not save on file.", err)
	}

	json.NewEncoder(w).Encode(weather)
}

func (server *WeatherServer) validateParams(r *http.Request) error {
	var errorMsg string
	if r.URL.Query().Get("latitude") == "" {
		errorMsg = "Latitude param not found."
	}
	if r.URL.Query().Get("longitude") == "" {
		errorMsg += "Longitude param not found."
	}
	if errorMsg != "" {
		return errors.New(errorMsg)
	}
	return nil
}

func (server *WeatherServer) makeHttpRequestWithContext(r *http.Request) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.open-meteo.com/v1/forecast", nil)
	if err != nil {
		log.Fatalf("Error while creating request with context: %v", err)
	}
	server.addParams(req.URL.Query(), r, req)
	return http.DefaultClient.Do(req)
}

func (server *WeatherServer) addParams(values url.Values, r *http.Request, req *http.Request) {
	values.Add("timezone", "America/Sao_Paulo")
	values.Add("current", "temperature_2m")
	values.Add("forecast_days", "1")
	values.Add("latitude", r.URL.Query().Get("latitude"))
	values.Add("longitude", r.URL.Query().Get("longitude"))
	req.URL.RawQuery = values.Encode()
}

func (server *WeatherServer) handleError(w http.ResponseWriter, httpStatus int, errorMsg string, err error) {
	w.WriteHeader(httpStatus)
	errorResponse := &errorResponse{
		ErrorMsg:       errorMsg,
		ErrorCause:     err.Error(),
		ErrorTimestamp: time.Now().Format(time.RFC850),
	}
	json.NewEncoder(w).Encode(errorResponse)
	file.RegisterLogFile(errorResponse, "ERROR: ", types.NewStandardLogFile())
}

func (server *WeatherServer) InitializeHandlers() {
	server.mux.HandleFunc("/weather", http.HandlerFunc(server.GetWeather))
}
