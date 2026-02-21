package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func Init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|'\n')
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|'\n')
}

func Info(msg string, args ...interface{}) {
	InfoLogger.Printf(msg, args...)
}

func Error(msg string, args ...interface{}) {
	ErrorLogger.Printf(msg, args...)
}
