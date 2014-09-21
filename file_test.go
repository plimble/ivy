package cayl

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTxtFile(t *testing.T) {
	assert := assert.New(t)

	fs := newFileStorage("./")

	str := "hello world"

	err := fs.Save("jack/test.txt", []byte(str))
	assert.NoError(err)

	data, err := fs.Load("jack/test.txt")
	assert.NoError(err)
	assert.Equal(data, []byte(str))

	err = fs.Delete("jack/test.txt")
	assert.NoError(err)

	os.RemoveAll("jack")

	fs = newFileStorage("./test2")

	err = fs.Save("jack/test.txt", []byte(str))
	assert.NoError(err)

	data, err = fs.Load("jack/test.txt")
	assert.NoError(err)
	assert.Equal(data, []byte(str))

	err = fs.Delete("jack/test.txt")
	assert.NoError(err)

	os.RemoveAll("test2")

}
