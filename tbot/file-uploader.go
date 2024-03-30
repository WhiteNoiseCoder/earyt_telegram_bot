package tbot

import (
	"io"
	"os"
)

type FileData struct {
	Path string
	Name string
}

func (fp FileData) NeedsUpload() bool {
	return true
}

func (fd FileData) UploadData() (string, io.Reader, error) {
	fileHandle, err := os.Open(fd.Path)
	if err != nil {
		return "", nil, err
	}

	name := fd.Name
	return name, fileHandle, err
}

func (fp FileData) SendData() string {
	panic("FilePath must be uploaded")
}
