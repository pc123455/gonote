package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	Fatal int = iota
	Error
	Warn
	Info
	Debug
	Trace
)

var logLevelMap = map[string]int{
	"fatal": Fatal,
	"error": Error,
	"warn":  Warn,
	"info":  Info,
	"debug": Debug,
	"Trace": Trace,
}

//var logger *log.Logger
var (
	level   int
	LogFile *os.File
)

func Initialize(filename string, logLevel string) {
	l, ok := logLevelMap[strings.ToLower(strings.TrimSpace(logLevel))]
	if !ok {
		panic("loglevel must be one of \"trace\", \"debug\", \"info\", \"warn\", \"error\",\"fatal\"")
	}
	level = l

	LogFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic("open logging file failed")
	}
	//defer func() {
	//	f.Close()
	//}()
	log.SetOutput(LogFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Close() {
	LogFile.Close()
}

func Tracef(format string, args ...interface{}) {
	if level >= Trace {
		log.SetPrefix("Trace ")
		log.Output(2, fmt.Sprintf(format, args))
	}
}

func Debugf(format string, args ...interface{}) {
	if level >= Debug {
		log.SetPrefix("Debug ")
		log.Output(2, fmt.Sprintf(format, args))
	}
}

func Infof(format string, args ...interface{}) {
	if level >= Info {
		log.SetPrefix("Info ")
		log.Output(2, fmt.Sprintf(format, args))
	}
}

func Warnf(format string, args ...interface{}) {
	if level >= Warn {
		log.SetPrefix("Warn ")
		log.Output(2, fmt.Sprintf(format, args))
	}
}

func Errorf(format string, args ...interface{}) {
	if level >= Error {
		log.SetPrefix("Error ")
		log.Output(2, fmt.Sprintf(format, args))
	}
}

func Fatalf(format string, args ...interface{}) {
	if level >= Error {
		log.SetPrefix("Fatal ")
		log.Output(2, fmt.Sprintf(format, args))
	}
}
