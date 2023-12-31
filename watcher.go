package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Router interface {
	Route(mes Message)
}
type WatchOption struct {
	recursive       bool
	processExisting bool
	ignoredDirs     map[string]struct{}
}

var DefaultWatcherOption = &WatchOption{
	recursive:       false,
	processExisting: false,
	ignoredDirs:     map[string]struct{}{".git": {}},
}

type DefaultWatcher struct {
	*fsnotify.Watcher
	router Router
}

func (watcher *DefaultWatcher) Watch(root string, options WatchOption) error {
	paths := watcher.getPaths(root, options)
	messages := make(chan Message)
	go watcher.watch(paths, messages)
	for {
		mes, ok := <-messages
		// The message channel was closed.
		if !ok {
			return nil
		}
		watcher.router.Route(mes)
	}
}

func (watcher *DefaultWatcher) getPaths(root string, options WatchOption) []string {
	if !options.recursive {
		return []string{root}
	}
	return walk(root, options.ignoredDirs)
}

func walk(root string, ignoreDirs map[string]struct{}) []string {
	paths := []string{}
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Something went wrong with the directory %v, permission perhaps ? The directory will not be watched", err)
			return err
		}
		if _, ok := ignoreDirs[path]; !d.IsDir() || ok {
			return err
		}
		paths = append(paths, path)
		return nil
	})
	return paths
}

func (watcher *DefaultWatcher) watch(paths []string, mes chan Message) {
	prettyPrint, err := PrettyPrint(paths)
	if err != nil {
		prettyPrint = fmt.Sprintf("%v", paths)
	}
	log.Printf("Listening to changes into paths\n%s", prettyPrint)

	defer watcher.Close()

	done := make(chan bool)

	// use goroutine to start the watcher
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				// The watcher is closed
				if !ok {
					close(mes)
					return
				}
				watcher.handleEvent(event, mes)
			case err, ok := <-watcher.Errors:
				// The watcher is closed
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
		log.Printf("The file %s is ignored because it's destination cannot be found", filename, filename)
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
	protocols := []string{"http", "ftp", "file"}
	index := rand.Intn(len(protocols))
	return protocols[index] + ":" + "/" + filename, nil
}

func PrettyPrint(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func NewWatcher(router Router) (*DefaultWatcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	watcher := DefaultWatcher{Watcher: fsWatcher, router: router}
	return &watcher, nil
}

func NewWatcherWithDefaultRouter(configPath string) (*DefaultWatcher, error) {
	content, err := os.ReadFile(configPath)

	if err != nil {
		return nil, err
	}

	configStr := string(content)

	config, err := ParseConfig(configStr)

	if err != nil {
		return nil, err
	}

	return NewWatcher(NewRouter(config))
}
