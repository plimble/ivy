package fileproxy

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadSource(t *testing.T) {
	fs := NewFileSystemSource("testsource")
	defer os.RemoveAll(fs.root)
	fileName := "test.txt"
	os.Mkdir(fs.root, 0755)
	err := ioutil.WriteFile(fs.root+"/"+fileName, []byte("TESTSOURCE"), 0755)
	assert.NoError(t, err)

	reader, err := fs.Load(fileName)
	assert.NotNil(t, reader)
	assert.NoError(t, err)
}

func TestLoadSourceNotExist(t *testing.T) {
	fs := NewFileSystemSource("testsource")
	fileName := "test.txt"
	reader, err := fs.Load(fileName)
	assert.Nil(t, reader)
	assert.Error(t, err)
}

func TestGetPath(t *testing.T) {
	fs := NewFileSystemSource("testsource")
	filename := fs.GetFilePath("test.txt")
	assert.Equal(t, "testsource/test.txt", filename)
}
