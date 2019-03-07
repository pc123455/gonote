package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	FatalLevel int = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

var logLevelMap = map[string]int{
	"fatal": FatalLevel,
	"error": ErrorLevel,
	"warn":  WarnLevel,
	"info":  InfoLevel,
	"debug": DebugLevel,
	"Trace": TraceLevel,
}

//var logger *log.Logger
var (
	level  int
	writer io.WriteCloser = os.Stdout
)

func SetOutput(out io.WriteCloser) {
	writer = out
	log.SetOutput(out)
}

func SetOutputFile(filename string) {
	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic("open logging file failed")
	}
	SetOutput(fd)
}

func SetLevel(logLevel string) {
	l, ok := logLevelMap[strings.ToLower(strings.TrimSpace(logLevel))]
	if !ok {
		panic("loglevel must be one of \"trace\", \"debug\", \"info\", \"warn\", \"error\",\"fatal\"")
	}
	level = l
}

func Init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Close() {
	writer.Close()
}

func Tracef(format string, args ...interface{}) {
	if level >= TraceLevel {
		log.SetPrefix("Trace ")
		log.Output(2, fmt.Sprintf(format, args))
	}
}

func Debugf(format string, args ...interface{}) {
	if level >= DebugLevel {
		log.SetPrefix("Debug ")
		log.Output(2, fmt.Sprintf(format, args))
	}
}

func Infof(format string, args ...interface{}) {
	if level >= InfoLevel {
		log.SetPrefix("Info ")
		log.Output(2, fmt.Sprintf(format, args))
	}
}

func Warnf(format string, args ...interface{}) {
	if level >= WarnLevel {
		log.SetPrefix("Warn ")
		log.Output(2, fmt.Sprintf(format, args))
	}
}

func Errorf(format string, args ...interface{}) {
	if level >= ErrorLevel {
		log.SetPrefix("Error ")
		log.Output(2, fmt.Sprintf(format, args))
	}
}

func Fatalf(format string, args ...interface{}) {
	log.SetPrefix("Fatal ")
	log.Output(2, fmt.Sprintf(format, args))
	os.Exit(-1)
}
