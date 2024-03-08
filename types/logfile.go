package types

import (
	"io"
	"log"
	"os"
)

type LogFile interface {
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	New(out io.Writer, prefix string, flag int) *log.Logger
	Println(log *log.Logger, v ...any)
}

type StandardLogFile struct{}

func NewStandardLogFile() *StandardLogFile {
	return &StandardLogFile{}
}

func (s *StandardLogFile) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (s *StandardLogFile) New(out io.Writer, prefix string, flag int) *log.Logger {
	return log.New(out, prefix, flag)
}

func (s *StandardLogFile) Println(log *log.Logger, v ...any) {
	log.Println(v...)
}
