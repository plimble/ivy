package fileproxy

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type FileSystemCache struct {
	root string
}

func NewFileSystemCache(root string) *FileSystemCache {
	if err := os.MkdirAll(root, 0755); err != nil {
		panic(err)
	}

	return &FileSystemCache{root}
}

func (fs *FileSystemCache) Save(bucket, filename, paramsStr string, data []byte) error {
	dir, filePath := filepath.Split(filename)
	filename = path.Join(fs.root, bucket, dir, paramsStr+filePath)
	dir = path.Join(fs.root, bucket, dir)

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

func (fs *FileSystemCache) Load(bucket, filename, paramsStr string) (*bytes.Buffer, error) {
	dir, filePath := filepath.Split(filename)
	filename = path.Join(fs.root, bucket, dir, paramsStr+filePath)

	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := bytes.NewBuffer(nil)
	buffer.ReadFrom(file)

	return buffer, nil
}

func (fs *FileSystemCache) Delete(bucket, filename string) error {
	filename = path.Join(fs.root, bucket, filename)
	ext := filepath.Ext(filename)
	if ext == "" {
		return fmt.Errorf("this is not file")
	}
	return os.Remove(filename)
}

func (fs *FileSystemCache) Flush() error {
	return os.RemoveAll(fs.root)
}
