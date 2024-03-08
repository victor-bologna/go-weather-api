package storage

import (
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/victor-bologna/go-weather-api/types"
)

func TestInit(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection.", err)
	}
	defer mockDB.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS WEATHER").WillReturnResult(sqlmock.NewResult(1, 1))

	mockPostgreStore := PostgreStore{DB: mockDB}

	if err = mockPostgreStore.Init(); err != nil {
		t.Errorf("Error while initiating mocked postgreDatabase: %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestInit_Error(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection.", err)
	}
	defer mockDB.Close()

	errorMsg := "Error when connecting postgre server."

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS WEATHER").WillReturnError(errors.New(errorMsg))

	mockPostgreStore := PostgreStore{DB: mockDB}

	if err = mockPostgreStore.Init(); err.Error() != errorMsg {
		t.Errorf("A different or no error occured when should occur: %s, but got: %s", errorMsg, err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestSaveWeather(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection.", err)
	}
	defer mockDB.Close()

	weatherMock := types.WeatherModel{
		Latitude:  -30.20,
		Longitude: -40.300,
		Timezone:  "America/Sao_Paulo",
		CurrentUnits: types.CurrentUnits{
			Time:          "iso8601",
			Interval:      "seconds",
			Temperature2M: "°C",
		},
		Current: types.Current{
			Time:          "2024-03-07T11:15",
			Interval:      900,
			Temperature2M: 25.5,
		},
	}

	mock.ExpectPrepare("INSERT INTO WEATHER ").ExpectExec().WithArgs(
		weatherMock.Latitude, weatherMock.Longitude, weatherMock.Timezone,
		weatherMock.CurrentUnits.Time, weatherMock.CurrentUnits.Interval, weatherMock.CurrentUnits.Temperature2M,
		weatherMock.Current.Time, weatherMock.Current.Interval, weatherMock.Current.Temperature2M).WillReturnResult(
		sqlmock.NewResult(1, 1))

	mockPostgreStore := PostgreStore{DB: mockDB}

	if err = mockPostgreStore.SaveWeather(&weatherMock); err != nil {
		t.Errorf("Error while saving weather in mocked postgreDatabase: %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestSaveWeather_ErrorPrepare(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection.", err)
	}
	defer mockDB.Close()

	errorMsg := "Error when preparing statement."

	mock.ExpectPrepare("INSERT INTO WEATHER ").WillReturnError(errors.New(errorMsg))

	mockPostgreStore := PostgreStore{DB: mockDB}

	if err = mockPostgreStore.SaveWeather(&types.WeatherModel{}); err.Error() != errorMsg {
		t.Errorf("A different or no error occured when should occur: %s, but got: %s", errorMsg, err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestSaveWeather_ErrorPrepareExec(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection.", err)
	}
	defer mockDB.Close()

	errorMsg := "sql: expected 9 arguments, got 8"

	mock.ExpectPrepare("INSERT INTO WEATHER ").ExpectExec().WillReturnError(errors.New(errorMsg))

	mockPostgreStore := PostgreStore{DB: mockDB}

	if err = mockPostgreStore.SaveWeather(&types.WeatherModel{}); err.Error() != errorMsg {
		t.Errorf("A different or no error occured when should occur: %s, but got: %s", errorMsg, err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
