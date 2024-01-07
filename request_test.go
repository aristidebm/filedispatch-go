package filedispatch

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	tmpDir, err := prepareTestDirTree("dir/to/walk/")

	if err != nil {
		fmt.Printf("unable to create test dir tree: (reason: %v)\n", err)
		return
	}

	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "example.txt")
	f, err := os.Create(filename)

	if err != nil {
		fmt.Printf("unable to create test file: (reason: %v)\n", err)
		return
	}

	print(io.ReadAll(f))

	// defer called are stacked, so this is called
	// before the former one and that is what we want
	defer f.Close()

	r, err := NewRequest(filename, "")

	if err != nil {
		fmt.Printf("unable to instanciate a new request (reason: %v)\n", err)
		return
	}

	assert.Equal(t, filename, r.Filename)
	assert.Equal(t, "", r.Destination)
	assert.Equal(t, int64(0), r.Size)
	assert.Equal(t, "text/plain", r.Type)
}
