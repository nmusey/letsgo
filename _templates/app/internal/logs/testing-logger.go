package logs

import "fmt"

type TestingLogger struct {}

func NewTestingLogger() Logger {
    return TestingLogger{}
}

func (tl TestingLogger) Log(message string, severity Severity) {
    fmt.Println(severity, " - ", message)
}
