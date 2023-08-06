package main

import (
	"log"
	"net/url"
	"strings"
)

type Worker interface {
	Process(mes Message)
	GetName() string
}

type DefaultRouter struct {
	workers map[string]Worker
}

func (router *DefaultRouter) Route(mes Message) {
	name, err := router.getProtocol(mes.destination)
	if err != nil {
		log.Printf("%s is an invalid path. Cannot route the message %v.", mes.destination, mes)
	}
	worker, ok := router.workers[name]
	if !ok {
		log.Printf("Unsupported protocol %s. Cannot route the message %v.", name, mes)
	}
	worker.Process(mes)
	log.Printf("The message %v is routed to %s worker", mes, name)
}

func (router *DefaultRouter) getProtocol(path string) (string, error) {
	parsedUrl, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	scheme := parsedUrl.Scheme
	if scheme == "" {
		scheme = strings.ToLower("file")
	}
	log.Printf("The parsed url is as follow %v", parsedUrl)
	return strings.ToLower(scheme), nil
}

func NewRouter() *DefaultRouter {
	workers := make(map[string]Worker)
	localWorker := LocalWorker{}
	workers[localWorker.GetName()] = &localWorker
	httpWorker := HttpWorker{}
	workers[httpWorker.GetName()] = &httpWorker
	ftpWorker := FtpWorker{}
	workers[ftpWorker.GetName()] = &ftpWorker
	router := &DefaultRouter{
		workers: workers,
	}
	return router
}
