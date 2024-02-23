package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Represents a logging instance.
type Logger struct {
	serviceName string
}

func CreateLogger(serviceName string) *Logger {
	return &Logger{serviceName}
}

func (l *Logger) Error(prefix string, err error) {
	file, line_no := getFileAndLineNo()
	fmt.Fprintf(os.Stderr, "%s %s %d: %s %v\n", l.getLogHeader(), file, line_no, prefix, err)
}

func (l *Logger) Message(message string) {
	fmt.Fprintf(os.Stdout, "%s %s\n", l.getLogHeader(), message)
}

func (l *Logger) getLogHeader() string {
	return "[" + l.serviceName + "]"
}

func getFileAndLineNo() (string, int) {
	_, file, line_no, _ := runtime.Caller(2)
	return filepath.Base(file), line_no
}
