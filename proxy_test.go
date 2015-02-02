package fileproxy_test

import (
	"bytes"
	"github.com/plimble/fileproxy"
	"github.com/stretchr/testify/assert"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func setup() *fileproxy.FileProxy {
	fsource := fileproxy.NewFileSystemSource("sourcefolder")
	csource := fileproxy.NewFileSystemCache("cachefolder")
	fconfig := &fileproxy.Config{
		IsDevelopment: false,
		HttpCache:     66000,
	}

	fp := fileproxy.New(
		fsource,
		csource,
		fconfig,
	)

	bucket := "bucket"
	os.MkdirAll(path.Join("sourcefolder", bucket), 0755)

	buffer := new(bytes.Buffer)
	png.Encode(buffer, image.NewRGBA(image.Rect(0, 0, 200, 200)))
	err := ioutil.WriteFile("sourcefolder/bucket/test.png", buffer.Bytes(), 0755)
	if err != nil {
		panic(err)
	}

	return fp
}

func teardown() {
	os.RemoveAll("sourcefolder")
	os.RemoveAll("cachefolder")
}

func TestProxyLoad(t *testing.T) {
	fp := setup()
	// defer teardown()

	bucket := "bucket"
	req, _ := http.NewRequest("GET", "bucket/test.png", nil)
	res := httptest.NewRecorder()
	fp.Get(bucket, "", "/test.png", res, req)

	assert.Equal(t, 200, res.Code)
}

func TestProxyResizeSize(t *testing.T) {
	bucket := "bucket"
	fp := setup()
	// defer teardown()

	req, _ := http.NewRequest("GET", "bucket/sourcetest/r_10x10/test.png", nil)
	res := httptest.NewRecorder()
	fp.Get(bucket, "r_10x10", "/test.png", res, req)

	assert.Equal(t, 200, res.Code)
}

func TestProxyLoadNotFound(t *testing.T) {
	bucket := "bucket"
	fp := setup()
	defer teardown()

	req, _ := http.NewRequest("GET", "/testnotfound.png", nil)
	res := httptest.NewRecorder()
	fp.Get(bucket, "", "/testnotfound.png", res, req)

	assert.Equal(t, 404, res.Code)
}
