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

	_, fsize, ftime, err := fs.Load(fileName)
	assert.NoError(t, err)
	assert.NotEmpty(t, fsize)
	assert.False(t, ftime.IsZero())
}

func TestLoadSourceNotExist(t *testing.T) {
	fs := NewFileSystemSource("testsource")
	fileName := "test.txt"
	_, fsize, ftime, err := fs.Load(fileName)
	assert.Error(t, err)
	assert.Empty(t, fsize)
	assert.True(t, ftime.IsZero())
}

func TestGetPath(t *testing.T) {
	fs := NewFileSystemSource("testsource")
	filename := fs.GetFilePath("test.txt")
	assert.Equal(t, "testsource/test.txt", filename)
}
