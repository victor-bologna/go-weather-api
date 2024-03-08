package main

import (
	"log"

	"github.com/victor-bologna/go-weather-api/api"
	"github.com/victor-bologna/go-weather-api/storage"
)

func main() {
	postgre, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal("Error when connecting postgre server.")
	}

	if err = postgre.Init(); err != nil {
		log.Fatal("Error when starting postgre server.")
	}

	server := api.NewWeatherAPI(":8080", postgre)
	server.StartServer()
}
