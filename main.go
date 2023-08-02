package main

import (
	//"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
)

type Message struct {
	filename    string
	destination string
}

type Router interface {
	route(mes Message) error
}

type DefautRouter struct {}

func (router *DefaultRouter) route(mes Message) {}

type Worker interface {
	entrust(mes Message)
	process() error
}

type WatchOption struct {
    recursive bool
    ignoredDirs [] string
}

var DefaulWatchOption = &WatchOption {
    recursive: false
    ignoreDirs: [".git"]
} 

type Watcher interface {
    watch(store sting, options WatchOption) error
}

type DefaultWatcher struct {
    router Router
}

func (watcher *DefaultWatcher) watch(store string, options WatchOption) error {}


func (worker *LocalStorageWorker) acquire(mes Message) {
	append(worker.state.queue, mes)
}


func main() {
	store := os.Args[1:]
    if len(store) < 1 {
        log.Fatal("You need to provide the store to watch")
    }
	watch(store[0])
}

func watch(store string) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(store)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
