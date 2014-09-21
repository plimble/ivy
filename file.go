package cayl

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileStorage struct {
	FilePath string
}

func newFileStorage(filePath string) *FileStorage {
	return &FileStorage{
		FilePath: filePath,
	}
}

func (fs *FileStorage) Save(filename string, data []byte) error {
	filename = fs.FilePath + "/" + filename
	exist, err := fs.Exist(filename)
	if err != nil {
		return err
	}
	if !exist {
		if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(filename, data, 0755)
}

func (fs *FileStorage) Load(filename string) (io.Reader, error) {
	filename = fs.FilePath + "/" + filename

	return os.Open(filename)
}

func (fs *FileStorage) Delete(filename string) error {
	filename = fs.FilePath + "/" + filename
	ext := filepath.Ext(filename)
	if ext == "" {
		return fmt.Errorf("this is not file")
	}
	return os.Remove(filename)
}

func (fs *FileStorage) DeleteFolder(dir string) error {
	return os.RemoveAll(dir)
}

func (fs *FileStorage) Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (fs *FileStorage) IsNotFound(err error) bool {
	return os.IsNotExist(err)
}
