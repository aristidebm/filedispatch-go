package main

import "log"

type Processor struct {
	queue []Message
}

func (processor *Processor) Process(mes Message) {
	log.Println("OK received I will process it.")
}

type LocalWorker struct {
	Processor
}

func (worker *LocalWorker) GetName() string { return "file" }

type HttpWorker struct {
	Processor
}

func (worker *HttpWorker) GetName() string { return "http" }

type FtpWorker struct {
	Processor
}

func (worker *FtpWorker) GetName() string { return "ftp" }
