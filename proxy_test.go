package ivy

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() *Ivy {
	source := newFakeSource()
	cache := newFakeCache()

	config := &Config{
		IsDevelopment: false,
		HttpCache:     66000,
	}

	iv := New(source, cache, config)

	return iv
}

func TestGetRaw(t *testing.T) {
	iv := setup()

	buffer := new(bytes.Buffer)
	png.Encode(buffer, image.NewRGBA(image.Rect(0, 0, 200, 200)))

	req, _ := http.NewRequest("GET", "bucket/_/test.png", nil)
	res := httptest.NewRecorder()

	iv.Get("bucket", "", "/test.png", res, req)

	assert.Equal(t, 200, res.Code)
}

func TestGetWithParams(t *testing.T) {
	iv := setup()

	req, _ := http.NewRequest("GET", "bucket/r_10x10/test.png", nil)
	res := httptest.NewRecorder()
	iv.Get("bucket", "r_10x10", "/test.png", res, req)

	assert.Equal(t, 200, res.Code)
}

func TestGetCacheNotFound(t *testing.T) {
	iv := setup()

	cache := newFakeCache()
	cache.err = ErrNotFound
	iv.Cache = cache

	source := newFakeSource()
	source.err = ErrNotFound
	iv.Source = source

	req, _ := http.NewRequest("GET", "bucket/r_100x100/testnotfound.png", nil)
	res := httptest.NewRecorder()
	iv.Get("bucket", "r_100x100", "/testnotfound.png", res, req)

	assert.Equal(t, 404, res.Code)
	assert.Equal(t, ErrNotFound.Error(), res.Body.String())
}

func TestGetSourceNotFound(t *testing.T) {
	iv := setup()

	cache := newFakeCache()
	cache.err = ErrNotFound
	iv.Cache = cache

	source := newFakeSource()
	source.err = ErrNotFound
	iv.Source = source

	req, _ := http.NewRequest("GET", "bucket/r_100x100/testnotfound.png", nil)
	res := httptest.NewRecorder()
	iv.Get("bucket", "r_100x100", "/testnotfound.png", res, req)

	assert.Equal(t, 404, res.Code)
	assert.Equal(t, ErrNotFound.Error(), res.Body.String())
}
