package fileproxy

import (
	"fmt"
	"io"
	"os"
	"path"
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
	filename = path.Join(fs.root, filename)
	dir := filepath.Dir(filename)

	_, err := os.Open(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL|os.O_SYNC, 0755)
	defer file.Close()
	if os.IsNotExist(err) || err == nil {
		if _, err := file.Write(data); err != nil {
			return err
		}
	}

	return err
}

func (fs *FileSystemCache) Load(filename string) (io.Reader, int64, time.Time, error) {
	filename = path.Join(fs.root, filename)

	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil, 0, time.Time{}, ErrNotFound
	}
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
	filename = path.Join(fs.root, filename)
	ext := filepath.Ext(filename)
	if ext == "" {
		return fmt.Errorf("this is not file")
	}
	return os.Remove(filename)
}

func (fs *FileSystemCache) Flush() error {
	return os.RemoveAll(fs.root)
}
