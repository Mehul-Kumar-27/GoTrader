package logger

import (
	"log"
	"os"
)

func CreateCustomLogger(info string) *log.Logger {
	logger := log.New(os.Stdout, info, log.LstdFlags|log.Lshortfile)
	return logger
}
