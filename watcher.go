package main

type WatchOption struct {
    recursive bool
    ignoredDirs [] string
}

var DefaulWatchOption = WatchOption {
    recursive: false,
    ignoreDirs: [".git"]
} 

type Watcher interface {
    watch(store sting, options WatchOption) error
}

type DefaultWatcher struct {
    router Router
}

func (watcher *DefaultWatcher) watch(store string, options WatchOption) error {
    for filename, destination := range defaultWatch {
        message := &Message {
            filename: filename,
            destination: destination
         }
         watcher.router.route(message)
  }
}
