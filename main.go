package main

import (
	"fmt"
	"log"
	"os"
)

type Message struct {
	filename    string
	destination string
}

func (mes Message) String() string {
	return fmt.Sprintf("%T{filename: %s, destination: %s}", mes, mes.filename, mes.destination)
}

type Watcher interface {
	Watch(store string, options WatchOption) error
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("You have to provide a directory to watch")
	}
	store := os.Args[1]
	watcher, err := NewWatcher()
	if err != nil {
		log.Fatal("Internal Error. Cannot watch any directory, retry later")
	}
	watcher.Watch(store, *DefaultWatcherOption)
}
