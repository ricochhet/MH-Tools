package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type LogLevel int

var LogCache []string

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Logger struct {
	MinLevel LogLevel
}

func NewLogger(minLevel LogLevel) *Logger {
	return &Logger{MinLevel: minLevel}
}

func (l *Logger) log(level LogLevel, message string, writer io.Writer) {
	if Stdout == nil {
		fmt.Fprintln(os.Stdout, "Logger STDOUT is nil")
		return
	}

	color := "\033[0m"
	levelName := "DEBUG"
	switch level {
	case DebugLevel:
		levelName = "DEBUG"
		color = "\033[0;34m"
	case InfoLevel:
		levelName = "INFO"
		color = "\033[0;32m"
	case WarnLevel:
		levelName = "WARN"
		color = "\033[0;33m"
	case ErrorLevel:
		levelName = "ERROR"
		color = "\033[0;31m"
	}

	if level >= l.MinLevel {
		oColor := fmt.Sprintf("[%s] %s[%s]%s %s\n", time.Now().Format("2006-01-02 15:04:05"), color, levelName, "\033[0m", message)
		oRaw := fmt.Sprintf("[%s] [%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), levelName, message)
		LogCache = append(LogCache, oRaw)
		fmt.Fprint(writer, oColor)
	}
}

func (l *Logger) Debug(message string) {
	l.log(DebugLevel, message, Stdout)
}

func (l *Logger) Info(message string) {
	l.log(InfoLevel, message, Stdout)
}

func (l *Logger) Warn(message string) {
	l.log(WarnLevel, message, Stdout)
}

func (l *Logger) Error(message string) {
	l.log(ErrorLevel, message, Stdout)
}

func (l *Logger) NewLine() {
	fmt.Fprint(Stdout, "\n")
}

func ClearCache() {
	LogCache = []string{}
}

var Stdout io.Writer
var SharedLogger = NewLogger(InfoLevel)
