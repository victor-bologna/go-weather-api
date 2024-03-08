package file

import (
	"io"
	"log"
	"os"

	"github.com/stretchr/testify/mock"
)

type mockLogFile struct {
	mock.Mock
}

type mockLogFileError struct {
	mock.Mock
}

func newMockLogFile() *mockLogFile { return &mockLogFile{} }

func (m *mockLogFile) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	args := m.Called(name, flag, perm)
	return args.Get(0).(*os.File), args.Error(1)
}

func (m *mockLogFile) New(out io.Writer, prefix string, flag int) *log.Logger {
	args := m.Called(out, prefix, flag)
	return args.Get(0).(*log.Logger)
}

func (m *mockLogFile) Println(log *log.Logger, v ...any) {
	m.Called(v)
}

func newMockLogFileError() *mockLogFileError { return &mockLogFileError{} }

func (m *mockLogFileError) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	args := m.Called(name, flag, perm)
	return nil, args.Error(1)
}

func (m *mockLogFileError) New(out io.Writer, prefix string, flag int) *log.Logger {
	panic("Unimplemented")
}

func (m *mockLogFileError) Println(log *log.Logger, v ...any) {
	panic("Unimplemented")
}
