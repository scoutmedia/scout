package logger

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

type Logger struct {
	Enviroment string
}

func NewLogger() *Logger {
	return &Logger{
		Enviroment: os.Getenv("env"),
	}
}

func (l *Logger) Init(dsn string) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      l.Enviroment,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(5 * time.Second)
}

func (l *Logger) Info(area string, msg string) {
	if l.Enviroment == "development" {
		log.Println(area, msg)
		return
	}
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelInfo)
		scope.SetTag("Scout", "scout-backend")
		scope.SetTag("Area", area)
		sentry.CaptureMessage(msg)
	})

}

func (l *Logger) Error(area string, msg string) {
	if l.Enviroment == "development" {
		log.Println(area, msg)
		return
	}
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelError)
		scope.SetTag("Scout", "scout-backend")
		scope.SetTag("Area", area)
		sentry.CaptureMessage(msg)
	})
}
