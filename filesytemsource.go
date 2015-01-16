package fileproxy

import (
	"io"
	"os"
	"time"
)

type FileSystemSource struct {
	root string
}

func NewFileSystemSource(root string) *FileSystemSource {
	return &FileSystemSource{root}
}

func (fs *FileSystemSource) Load(filename string) (io.Reader, int64, time.Time, error) {
	filename = fs.root + "/" + filename

	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, time.Time{}, err
	}
	ft, err := file.Stat()
	if err != nil {
		return nil, 0, time.Time{}, err
	}

	return file, ft.Size(), ft.ModTime(), nil
}

func (fs *FileSystemSource) GetFilePath(filename string) string {
	return fs.root + "/" + filename
}
