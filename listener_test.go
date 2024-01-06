package filedispatch

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeRouter struct{}

func (r FakeRouter) Route(path string) error {
	return nil
}

func TestIgnoreDirs(t *testing.T) {
	router := FakeRouter{}
	listener, err := NewListener(router)

	if err != nil {
		log.Printf("unable to instiante the listener")
		return
	}

	listener = listener.WithIgnoreDirs("/tmp/dir1", "/tmp/dir3").
		WithIgnoreDirs("/tmp/dir4", "dir5", "dir2/dir8").WithIgnoreDirs("^private$")

	assert.True(t, listener.canWatchDir("/tmp"))
	assert.True(t, listener.canWatchDir("/mnt"))
	assert.True(t, listener.canWatchDir("dir7"))
	assert.True(t, listener.canWatchDir("/tmp/storage/dir1"))
	assert.True(t, listener.canWatchDir("/tmp/private"))

	assert.False(t, listener.canWatchDir("/tmp/dir1"))
	assert.False(t, listener.canWatchDir("dir2/dir8"))
	assert.False(t, listener.canWatchDir("/tmp/dir2/dir8"))
	assert.False(t, listener.canWatchDir("/home/dir2/dir8"))
	assert.False(t, listener.canWatchDir("dir5"))
	assert.False(t, listener.canWatchDir("/home/user/dir5"))
	assert.False(t, listener.canWatchDir("/tmp/dir5"))
	assert.False(t, listener.canWatchDir("private"))
}

func TestIgnoreFiles(t *testing.T) {
	router := FakeRouter{}
	listener, err := NewListener(router)

	if err != nil {
		log.Printf("unable to instiante the listener")
		return
	}

	listener = listener.WithIgnoreFiles("/tmp/file1.txt", "/tmp/file3.txt").
		WithIgnoreFiles("/tmp/file4.txt", "file5.txt", "dir2/file8.txt").WithIgnoreFiles("^private.txt$")

	assert.True(t, listener.canHandleFile("/tmp/storage/file1.txt"))
	assert.True(t, listener.canHandleFile("/tmp/private.txt"))

	assert.False(t, listener.canHandleFile("/tmp/file1.txt"))
	assert.False(t, listener.canHandleFile("dir2/file8.txt"))
	assert.False(t, listener.canHandleFile("/tmp/dir2/file8.txt"))
	assert.False(t, listener.canHandleFile("/home/dir2/file8.txt"))
	assert.False(t, listener.canHandleFile("file5.txt"))
	assert.False(t, listener.canHandleFile("/home/user/file5.txt"))
	assert.False(t, listener.canHandleFile("/tmp/file5.txt"))
	assert.False(t, listener.canHandleFile("private.txt"))
}

func TestRecusrvice(t *testing.T) {
	tmpDir, err := prepareTestDirTree("dir/to/walk/")

	if err != nil {
		fmt.Printf("unable to create test dir tree: %v\n", err)
		return
	}

	defer os.RemoveAll(tmpDir)

	router := FakeRouter{}
	listener, err := NewListener(router)

	if err != nil {
		log.Printf("unable to instiante the listener")
		return
	}

	listener = listener.WithRecursive(false)
	paths, err := listener.getPaths(tmpDir)

	if err != nil {
		log.Printf("unable paths")
		return
	}

	assert.Len(t, paths, 1)
	assert.Equal(t, paths[0], tmpDir)

	listener = listener.WithRecursive(true)
	paths, err = listener.getPaths(tmpDir)

	if err != nil {
		log.Printf("unable to get the paths")
		return
	}
	expected := []string{tmpDir, filepath.Join(tmpDir, "dir"), filepath.Join(tmpDir, "dir/to"), filepath.Join(tmpDir, "dir/to/walk")}
	assert.Len(t, paths, 4)
	assert.Equal(t, paths, expected)
}

func TestRecusrviceAndIgnoreDir(t *testing.T) {
	//
}

func TestRecusrviceAndIgnoreFile(t *testing.T) {
	//
}

// source https://github.com/golang/go/blob/master/src/path/filepath/example_unix_walk_test.go#L16C1-L29C2
func prepareTestDirTree(tree string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", fmt.Errorf("error creating temp directory: %v\n", err)
	}

	err = os.MkdirAll(filepath.Join(tmpDir, tree), 0755)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}
	return tmpDir, nil
}
