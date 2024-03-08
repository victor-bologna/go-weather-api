package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/victor-bologna/go-weather-api/types"
)

type PostgreStore struct {
	DB *sql.DB
}

func NewPostgresStore() (*PostgreStore, error) {
	connStr := "user=root dbname=weather password=weather_pass sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgreStore{
		DB: db,
	}, nil
}

func (postgres *PostgreStore) Init() error {
	return postgres.createWeatherTable()
}

func (postgres *PostgreStore) createWeatherTable() error {
	createSQL := `CREATE TABLE IF NOT EXISTS WEATHER (
		id serial primary key,
		"latitude" double precision,
		"longitude" double precision,
		"timezone" text,
		"current_units.time" varchar(50),
		"current_units.interval" varchar(50),
		"current_units.temperature_2m" varchar(50),
		"current.time" varchar(50),
		"current.interval" bigint,
		"current.temperature_2m" double precision
	  );`
	_, err := postgres.DB.Exec(createSQL)
	return err
}

func (postgres *PostgreStore) SaveWeather(weather *types.WeatherModel) error {
	stmt, err := postgres.DB.Prepare(`
        INSERT INTO WEATHER (
            latitude, longitude, timezone,
            "current_units.time", "current_units.interval", "current_units.temperature_2m",
            "current.time", "current.interval", "current.temperature_2m"
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(weather.Latitude, weather.Longitude, weather.Timezone,
		weather.CurrentUnits.Time, weather.CurrentUnits.Interval, weather.CurrentUnits.Temperature2M,
		weather.Current.Time, weather.Current.Interval, weather.Current.Temperature2M)
	return err
}
