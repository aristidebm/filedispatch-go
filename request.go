package filedispatch

import (
	"fmt"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

type Request struct {
	Filename    string
	Destination string
	Type        string
	Size        int64
}

func NewRequest(filename string, destination string) (Request, error) {
	size, err := getSize(filename)

	if err != nil {
		fmt.Printf("I am here One %v \n", err)
		return Request{}, err
	}

	fType, err := getType(filename)

	if err != nil {
		fmt.Printf("I am here Two %v \n", err)
		return Request{}, err
	}

	return Request{
		Filename:    filename,
		Destination: destination,
		Type:        fType,
		Size:        size,
	}, nil
}

func getType(filename string) (string, error) {
	mtype, err := mimetype.DetectFile(filename)
	if err != nil {
		return "", fmt.Errorf("unable to get the type of the file %s (reason: %v)", filename, err)
	}
	// fmt.Println(mtype.Is("application/octet-stream"), mtype.String(), os.IsNotExist(err))
	return mtype.String(), nil
}

func getSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, fmt.Errorf("unable to get the size of the file %s (reason: %v)", filename, err)
	}
	return info.Size(), nil
}
