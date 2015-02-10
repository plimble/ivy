package ivy

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestFileSourceLoad(t *testing.T) {
	bucket := "bucket"
	fs := NewFileSystemSource("testsource")
	defer os.RemoveAll(fs.root)
	fileName := "test.txt"
	err := os.MkdirAll(path.Join(fs.root, bucket), 0755)
	assert.NoError(t, err)

	err = ioutil.WriteFile(path.Join(fs.root, bucket, fileName), []byte("TESTSOURCE"), 0755)
	assert.NoError(t, err)

	reader, err := fs.Load(bucket, fileName)
	assert.NotNil(t, reader)
	assert.NoError(t, err)
}

func TestFileSourceLoadNotExist(t *testing.T) {
	bucket := "bucket"
	fs := NewFileSystemSource("testsource")
	fileName := "test.txt"
	reader, err := fs.Load(bucket, fileName)
	assert.Nil(t, reader)
	assert.Error(t, err)
}

func TestFileSourceGetPath(t *testing.T) {
	bucket := "bucket"
	fs := NewFileSystemSource("testsource")
	filename := fs.GetFilePath(bucket, "test.txt")
	assert.Equal(t, "testsource/bucket/test.txt", filename)
}
