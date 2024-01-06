package filedispatch

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type Router interface {
	Route(path string) error
}

type listenerOption struct {
	recursive   bool
	ignoreDirs  []*regexp.Regexp
	ignoreFiles []*regexp.Regexp
}

func newListenerOption() (*listenerOption, error) {
	return &listenerOption{
		recursive: false,
	}, nil
}

type Listener struct {
	option *listenerOption
	mux    Router
}

func NewListener(router Router) (Listener, error) {
	option, err := newListenerOption()

	if err != nil {
		return Listener{}, fmt.Errorf("cannot instanciate the listener (reason: %s)", err.Error())
	}

	return Listener{mux: router, option: option}, nil
}

func (l Listener) getPaths(path string) ([]string, error) {
	if !l.option.recursive {
		return []string{path}, nil
	}
	return l.walk(path)
}

func (l Listener) walk(root string) ([]string, error) {
	paths := []string{}

	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("cannot listen to subdirectories of the directory %s (reason: %s)", root, err.Error())
		}

		// We are only interested in Dirs
		if !info.IsDir() {
			return nil
		}

		if l.canWatchDir(path) {
			paths = append(paths, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}

func (l Listener) canWatchDir(path string) bool {
	for _, p := range l.option.ignoreDirs {
		if p.MatchString(path) {
			return false
		}
	}
	return true
}

func (l Listener) canHandleFile(path string) bool {
	for _, p := range l.option.ignoreFiles {
		if p.MatchString(path) {
			return false
		}
	}
	return true
}

func (l Listener) watch(paths []string) error {

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return fmt.Errorf("cannot listen to directories %v (reason: %s)", paths, err.Error())
	}

	defer watcher.Close()

	log.Printf("Listerning to changes into directories\n %v", paths)

	for _, p := range paths {
		if err := watcher.Add(p); err != nil {
			return fmt.Errorf("cannot listen the directory %s (reason: %s)", p, err.Error())
		}
	}

	var wg sync.WaitGroup

	// use a goroutine to listen to changes
	// and process them
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case event, ok := <-watcher.Events:

				// the events channel is closed, we cannot receive
				// the event anymore, end the goroutine
				if !ok {
					return
				}

				if event.Has(fsnotify.Create) {
					go func() {
						if err := l.mux.Route(event.Name); err != nil {
							watcher.Errors <- err
						}
					}()
				}

			case err, ok := <-watcher.Errors:
				// the error channel is closed, we cannot receive
				// the errors anymore, end the goroutine
				if !ok {
					return
				}
				log.Print(err)
			}
		}
	}()

	wg.Wait()

	return nil
}

func (l Listener) Listen(path string) error {
	paths, err := l.getPaths(path)

	if err != nil {
		return err
	}

	if err = l.watch(paths); err != nil {
		return err
	}

	return nil
}

func (l Listener) WithRecursive(value bool) Listener {
	l.option.recursive = value
	return l
}

func (l Listener) WithIgnoreDirs(value ...string) Listener {
	for _, v := range value {
		p, err := regexp.Compile(v)

		if err != nil {
			log.Printf("cannot ignore dir %s (reason: %s)", v, err.Error())
			continue
		}
		l.option.ignoreDirs = append(l.option.ignoreDirs, p)
	}
	return l
}

func (l Listener) WithIgnoreFiles(value ...string) Listener {
	for _, v := range value {
		p, err := regexp.Compile(v)

		if err != nil {
			log.Printf("cannot ignore file %s (reason: %s)", v, err.Error())
			continue
		}

		l.option.ignoreFiles = append(l.option.ignoreFiles, p)
	}
	return l
}
