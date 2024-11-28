package logs

import (
	"log"
	"time"
)

type Logger interface {
    Log(string, Severity)
}

type DefaultLogger struct {}

func NewDefaultLogger() DefaultLogger {
    return DefaultLogger{}
}

type Severity int
const (
    Debug = iota
    Warning
    Error
    Alert
)

func (l DefaultLogger) Log(message string, severity Severity) {
    log.Println(time.Now(), " - ", severity, " - ", message)
}
