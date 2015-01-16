package fileproxy

// import (
// 	"bytes"
// 	"errors"
// 	"github.com/stretchr/testify/assert"
// 	"image"
// 	"image/png"
// 	"io"
// 	"testing"
// )

// type fakeStorage struct{}

// func (fs *fakeStorage) Save(filename string, data []byte) error {
// 	return nil
// }

// func (fs *fakeStorage) Load(filename string) (io.Reader, error) {
// 	params, _ := ParseParams("w_100,h_100")
// 	if filename == sourcePath("file1.png") || filename == cachePath(params, "file2.png", ".png") {
// 		buffer := bytes.NewBuffer(nil)
// 		png.Encode(buffer, image.NewRGBA(image.Rect(0, 0, 200, 200)))
// 		return buffer, nil
// 	}

// 	return nil, NotFound
// }

// func (fs *fakeStorage) Delete(filename string) error {
// 	if filename == sourcePath("file.jpg") {
// 		return nil
// 	} else if filename == sourcePath("file") {
// 		return errors.New("not a file")
// 	} else if filename == sourcePath("file2.jpg") {
// 		return NotFound
// 	}
// 	return nil
// }

// func (fs *fakeStorage) DeleteFolder(dir string) error {
// 	return nil
// }

// func (fs *fakeStorage) Exist(path string) (bool, error) {
// 	if path == sourcePath("test/test.jpg") {
// 		return true, nil
// 	}

// 	return false, nil
// }

// func (fs *fakeStorage) IsNotFound(err error) bool {
// 	if err != NotFound {
// 		return false
// 	}

// 	return true
// }

// func TestSave(t *testing.T) {
// 	assert := assert.New(t)
// 	p := New(Config{})
// 	p.File = &fakeStorage{}

// 	buffer := bytes.NewBuffer(nil)
// 	png.Encode(buffer, image.NewRGBA(image.Rect(0, 0, 200, 200)))

// 	//save success
// 	result, err := p.Save(File, "test.jpg", buffer)
// 	assert.NoError(err)
// 	assert.NotNil(result)

// 	//is already exist
// 	result, err = p.Save(File, "test/test.jpg", buffer)
// 	assert.Error(err)
// 	assert.Nil(result)
// }

// func TestLoad(t *testing.T) {
// 	assert := assert.New(t)
// 	p := New(Config{})
// 	p.File = &fakeStorage{}

// 	// load from source
// 	result, err := p.Load(File, "_", "file1.png", ".png")
// 	assert.NoError(err)
// 	assert.NotNil(result)

// 	// load from cache
// 	result, err = p.Load(File, "w_100,h_100", "file2.png", ".png")
// 	assert.NoError(err)
// 	assert.NotNil(result)

// 	//load from source and process
// 	result, err = p.Load(File, "w_100,h_100", "file1.png", ".png")
// 	assert.NoError(err)
// 	assert.NotNil(result)

// 	//not found
// 	result, err = p.Load(File, "w_100,h_100", "file3.png", ".png")
// 	assert.Error(err)
// 	assert.Nil(result)

// 	//parse error
// 	result, err = p.Load(File, "w_100h_100", "file1.png", ".png")
// 	assert.Error(err)
// 	assert.Nil(result)
// }

// func TestDelete(t *testing.T) {
// 	assert := assert.New(t)
// 	p := New(Config{})
// 	p.File = &fakeStorage{}

// 	//delete success
// 	err := p.Delete(File, "file.jpg")
// 	assert.NoError(err)

// 	//not found
// 	err = p.Delete(File, "file2.jpg")
// 	assert.Error(err)
// 	assert.Equal(err, NotFound)

// 	//not file
// 	err = p.Delete(File, "file")
// 	assert.Error(err)
// }
