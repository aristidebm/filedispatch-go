package main

import (
	"log"
	"strings"
)

type MessageLogger interface {
	LogMessage(mes Message) error
	GetName() string
}

type Processor struct {
	queue   []Message
	loggers map[string]MessageLogger
}

func (processor *Processor) ProcessMessage(mes Message) {
	log.Println("OK received I will process it.")
}

func (processor *Processor) initLoggers() {
	// TODO: How can we initLoggers without relying on a factory method or calling it from ProcessMessage
	processor.loggers = map[string]MessageLogger{}
	loggers := []MessageLogger{&JsonMessageLogger{}, &CsvMessageLogger{}, &ConsoleMessageLogger{}, &SqliteMessageLogger{}}
	for _, logger := range loggers {
		processor.registerLogger(logger)
	}
}

func (processor *Processor) registerLogger(logger MessageLogger) {
	processor.loggers[logger.GetName()] = logger
}

type LocalWorker struct {
	Processor
}

func (worker *LocalWorker) GetName() string { return strings.ToLower("file") }

type HttpWorker struct {
	Processor
}

func (worker *HttpWorker) GetName() string { return strings.ToLower("http") }

type FtpWorker struct {
	Processor
}

func (worker *FtpWorker) GetName() string { return strings.ToLower("ftp") }
