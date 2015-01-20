package fileproxy

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestSave(t *testing.T) {
	fs := NewFileSystemCache("testcache")
	pathFile := "subfolder/test.txt"
	paramsStr := ""
	dir, filePath := filepath.Split(pathFile)
	fullPath := path.Join(fs.root, dir, paramsStr+filePath)
	defer os.RemoveAll("testcache")

	err := fs.Save(pathFile, "", []byte("TESTCACHE"))
	assert.NoError(t, err)

	byteString, err := ioutil.ReadFile(fullPath)
	assert.NoError(t, err)
	assert.Equal(t, "TESTCACHE", string(byteString))
}

func TestSaveExist(t *testing.T) {
	fs := NewFileSystemCache("testcache")
	defer os.RemoveAll("testcache")

	err := fs.Save("test.txt", "", []byte("TESTCACHE"))
	err = fs.Save("test.txt", "", []byte("TESTCACHE"))
	assert.Error(t, err)
}

func TestLoadCache(t *testing.T) {
	fs := NewFileSystemCache("testcache")
	pathFile := "subfolder/test.txt"
	defer os.RemoveAll("testcache")

	err := fs.Save(pathFile, "", []byte("TESTCACHE"))
	assert.NoError(t, err)

	reader, err := fs.Load(pathFile, "")

	assert.NoError(t, err)

	data := bytes.NewBuffer(nil)
	data.ReadFrom(reader)
	assert.Equal(t, "TESTCACHE", data.String())
}

func TestLoadNotExist(t *testing.T) {
	fs := NewFileSystemCache("testcache")
	reader, err := fs.Load("test.txt", "")
	assert.Error(t, err)
	assert.Nil(t, reader)
}

func TestFlush(t *testing.T) {
	fs := NewFileSystemCache("testcache")

	err := fs.Save("test.txt", "", []byte("TESTCACHE"))
	assert.NoError(t, err)

	err = fs.Save("test2.txt", "", []byte("TESTCACHE"))
	assert.NoError(t, err)

	err = fs.Flush()
	assert.NoError(t, err)

	_, err = os.Open(fs.root)
	assert.True(t, os.IsNotExist(err))
}
