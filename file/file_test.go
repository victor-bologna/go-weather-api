package file

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterLogFile_Success(t *testing.T) {
	weatherMock := "mocked weather data"
	status := "INFO: "

	var mockFile = &os.File{}
	var mockLogger = &log.Logger{}

	mockLogFile := NewMockLogFile()
	mockLogFile.On("OpenFile", "weather.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(mockFile, nil)
	mockLogFile.On("New", mockFile, status, log.Ldate|log.Ltime|log.Lmicroseconds).Return(mockLogger)
	mockLogFile.On("Println", mock.Anything)

	err := RegisterLogFile(weatherMock, status, mockLogFile)

	assert.NoError(t, err)
	mockLogFile.AssertExpectations(t)
	mockLogFile.AssertNumberOfCalls(t, "OpenFile", 1)
	mockLogFile.AssertNumberOfCalls(t, "New", 1)
	mockLogFile.AssertNumberOfCalls(t, "Println", 1)
}

func TestRegisterLogFile_Error(t *testing.T) {
	weatherMock := "mocked weather data"
	status := "Error: "

	err2 := errors.New("open weather.log: The system cannot find the file specified.")

	mockLogFile := NewMockLogFileError()
	mockLogFile.On("OpenFile", "weather.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(nil, err2)

	err := RegisterLogFile(weatherMock, status, mockLogFile)

	assert.Equal(t, err.Error(), err2.Error())
	mockLogFile.AssertExpectations(t)
	mockLogFile.AssertNumberOfCalls(t, "OpenFile", 1)
	mockLogFile.AssertNumberOfCalls(t, "New", 0)
	mockLogFile.AssertNumberOfCalls(t, "Println", 0)
}
