package main

import (
     "log"
     fd "example.com/project/filedispatch"
)

func main() {
    fd.NewListener(NewFakeMux())
	mux := NewFakeMux()
	listener, err := fd.NewListener(mux)

	if err != nil {
		log.Fatal(err)
	}

	listener = listener.WithRecursive(true).WithIgnoreDirs(
		"dir2",
		"dir3",
	)
	if err = listener.Listen("/tmp/storage"); err != nil {
		log.Fatal(err)
	}
}

type FakeMux struct{}

func (m FakeMux) Route(path string) error {
    log.Printf("Routing %s", path)
    return nil
}

func NewFakeMux() *FakeMux {
    return &FakeMux{}
}
