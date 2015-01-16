package fileproxy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileSystemCache struct {
	root string
}

func NewFileSystemCache(root string) *FileSystemCache {
	return &FileSystemCache{root}
}

func (fs *FileSystemCache) Save(filename string, data []byte) error {
	filename = fs.root + "/" + filename

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.Stat()
	if os.IsExist(err) {
		_, err = file.Write(data)
		return err
	}

	return err
}

func (fs *FileSystemCache) Load(filename string) (io.Reader, int64, time.Time, error) {
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

func (fs *FileSystemCache) Delete(filename string) error {
	filename = fs.root + "/" + filename
	ext := filepath.Ext(filename)
	if ext == "" {
		return fmt.Errorf("this is not file")
	}
	return os.Remove(filename)
}

func (fs *FileSystemCache) Flush() error {
	return os.RemoveAll(fs.root)
}
