package fileproxy

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	accessKey = os.Getenv("AWS_ACCESS_KEY")
	secretKey = os.Getenv("oUeb1WiAykPGE0HoOVvMsu24u1Q")
)

func TestS3SourceLoad(t *testing.T) {
	bucket := "plimble"
	fs := NewS3Source(accessKey, secretKey)

	reader, err := fs.Load(bucket, "test/a123/v1/test.jpg")
	assert.NotNil(t, reader)
	assert.NoError(t, err)
}

func TestS3SourceLoadNotExist(t *testing.T) {
	bucket := "plimble"
	fs := NewS3Source(accessKey, secretKey)

	reader, err := fs.Load(bucket, "test/a123/v1/test2.jpg")
	assert.Nil(t, reader)
	assert.Error(t, err)
}

func TestS3SourceGetPath(t *testing.T) {
	bucket := "bucket"
	fs := NewS3Source(accessKey, secretKey)
	filename := fs.GetFilePath(bucket, "test/a123/v1/test.jpg")
	assert.Equal(t, filename, "test/a123/v1/test.jpg")
}
