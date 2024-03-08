package types

type CurrentUnits struct {
	Time          string `json:"time"`
	Interval      string `json:"interval"`
	Temperature2M string `json:"temperature_2m"`
}

type Current struct {
	Time          string  `json:"time"`
	Interval      int     `json:"interval"`
	Temperature2M float64 `json:"temperature_2m"`
}

type WeatherModel struct {
	Latitude     float64      `json:"latitude"`
	Longitude    float64      `json:"longitude"`
	Timezone     string       `json:"timezone"`
	CurrentUnits CurrentUnits `json:"current_units"`
	Current      Current      `json:"current"`
}
