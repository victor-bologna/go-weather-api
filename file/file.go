package file

import (
	"encoding/json"
	"log"
	"os"

	"github.com/victor-bologna/go-weather-api/types"
)

func RegisterLogFile(body interface{}, status string, l types.LogFile) error {
	file, err := l.OpenFile("weather.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer file.Close()
	data, _ := json.Marshal(body)
	logger := l.New(file, status, log.Ldate|log.Ltime|log.Lmicroseconds)
	l.Println(logger, string(data))
	return nil
}
