package log

import "fmt"

const (
	DEBUG = iota
	INFO
	ERROR
)

var Level int = INFO

func SetLogLevel(level string) {
	switch level {
	case "DEBUG":
		Level = DEBUG
	case "INFO":
		Level = INFO
	case "ERROR":
		Level = ERROR
	}
}

func Debug(format string, a ...interface{}) {
	if Level <= DEBUG {
		fmt.Printf(format, a...)
	}
}

func Error(format string, a ...interface{}) {
	if Level <= ERROR {
		fmt.Printf(format, a...)
	}
}
