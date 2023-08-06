package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type Router interface {
	Route(mes Message)
}
type WatchOption struct {
	recursive   bool
	ignoredDirs []string
}

var DefaultWatcherOption = &WatchOption{
	recursive:   false,
	ignoredDirs: []string{".git"},
}

type DefaultWatcher struct {
	*fsnotify.Watcher
	router Router
}

func (watcher *DefaultWatcher) Watch(store string, options WatchOption) error {
	paths := watcher.getPaths(store, options)
	messages := make(chan Message)
	go watcher.watch(paths, messages)
	for mes := range messages {
		watcher.router.Route(mes)
	}
	return nil
}

func (watcher *DefaultWatcher) getPaths(store string, options WatchOption) []string {
	paths := []string{store}
	return paths
}

func (watcher *DefaultWatcher) watch(paths []string, mes chan Message) {
	log.Printf("Listening to paths %v", paths)

	defer watcher.Close()

	done := make(chan bool)

	// use goroutine to start the watcher
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok { // The watcher is closed
					return
				}
				watcher.handleEvent(event, mes)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				watcher.handleError(err)
			}
		}
	}()
	watcher.addPaths(paths)
	<-done
}

func (watcher *DefaultWatcher) addPaths(paths []string) {
	for _, path := range paths {
		if err := watcher.Add(path); err != nil {
			log.Fatal(err)
		}
	}
}

func (watcher *DefaultWatcher) handleEvent(event fsnotify.Event, mes chan Message) {

	if !event.Has(fsnotify.Create) {
		return
	}

	filename := event.Name
	destination, err := watcher.getDestination(filename)

	if err != nil {
		log.Printf("Cannot find the destination of the file %s. The file %s is ignored", filename, filename)
		return
	}

	mes <- Message{
		filename:    filename,
		destination: destination,
	}
}

func (watcher *DefaultWatcher) handleError(err error) {
	log.Fatal("There is an error")
}

func (watcher *DefaultWatcher) getDestination(filename string) (string, error) {
	return filename, nil
}

func NewWatcher(routers ...Router) (*DefaultWatcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	if len(routers) > 0 {
		return &DefaultWatcher{
			Watcher: fsWatcher,
			router:  routers[0],
		}, nil
	}

	return &DefaultWatcher{
		Watcher: fsWatcher,
		router:  NewRouter(),
	}, nil
}
