package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIgnoreDirs(t *testing.T) {
	router := NewFakeMux()
	listener, err := NewListener(router)

	if err != nil {
		log.Fatal(err)
	}

	listener = listener.WithIgnoreDirs("/tmp/dir1", "/tmp/dir3").
		WithIgnoreDirs("/tmp/dir4", "dir5", "dir2/dir8")

	assert.True(t, listener.canWatchDir("/tmp"))
	assert.True(t, listener.canWatchDir("/mnt"))
	assert.True(t, listener.canWatchDir("dir7"))
	assert.True(t, listener.canWatchDir("/tmp/storage/dir1"))

	assert.False(t, listener.canWatchDir("/tmp/dir1"))
	assert.False(t, listener.canWatchDir("dir2/dir8"))
	assert.False(t, listener.canWatchDir("/tmp/dir2/dir8"))
	assert.False(t, listener.canWatchDir("/home/dir2/dir8"))
	assert.False(t, listener.canWatchDir("dir5"))
	assert.False(t, listener.canWatchDir("/home/user/dir5"))
	assert.False(t, listener.canWatchDir("/tmp/dir5"))
}

func TestIgnoreFiles(t *testing.T) {

}

func TestIgnoreRecusrvice(t *testing.T) {

}

func TestIgnoreWatch(t *testing.T) {

}
