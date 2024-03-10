package file

import (
	"io"
	"log"
	"os"

	"github.com/stretchr/testify/mock"
)

type MockLogFile struct {
	mock.Mock
}

type MockLogFileError struct {
	mock.Mock
}

func NewMockLogFile() *MockLogFile { return &MockLogFile{} }

func (m *MockLogFile) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	args := m.Called(name, flag, perm)
	return args.Get(0).(*os.File), args.Error(1)
}

func (m *MockLogFile) New(out io.Writer, prefix string, flag int) *log.Logger {
	args := m.Called(out, prefix, flag)
	return args.Get(0).(*log.Logger)
}

func (m *MockLogFile) Println(log *log.Logger, v ...any) {
	m.Called(v)
}

func NewMockLogFileError() *MockLogFileError { return &MockLogFileError{} }

func (m *MockLogFileError) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	args := m.Called(name, flag, perm)
	return nil, args.Error(1)
}

func (m *MockLogFileError) New(out io.Writer, prefix string, flag int) *log.Logger {
	panic("Unimplemented")
}

func (m *MockLogFileError) Println(log *log.Logger, v ...any) {
	panic("Unimplemented")
}
