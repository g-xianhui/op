package main

import (
	oslog "log"
	"os"
)

const (
	_ = iota
	DEBUG
	INFO
	ERROR
)

var levels = map[string]int{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"ERROR": ERROR,
}

var logLevel int
var logger *oslog.Logger

func logInit(file string, level string) {
	logfile, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		oslog.Fatal(err)
	}
	logger = oslog.New(logfile, "", oslog.Ldate|oslog.Ltime|oslog.Llongfile)
	logLevel = levels[level]
}

func log(level int, format string, args ...interface{}) {
	if level < logLevel {
		return
	}
	logger.Printf(format, args...)
}
