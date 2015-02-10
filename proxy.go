package ivy

import (
	"bytes"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

//ErrNotFound is error not found
var ErrNotFound = errors.New("not found")

//Source interface
type Source interface {
	Load(bucket string, filename string) (*bytes.Buffer, error)
	GetFilePath(bucket string, filename string) string
}

//Cache interface
type Cache interface {
	Save(bucket, filename, paramsStr string, file []byte) error
	Load(bucket, filename, paramsStr string) (*bytes.Buffer, error)
	Delete(bucket, filename string) error
	Flush(bucket string) error
}

//Ivy main struct of this package
type Ivy struct {
	Source Source
	Cache  Cache
	Config *Config
}

//Config config for ivy
type Config struct {
	HTTPCache     int64 //Enable http cache in second
	IsDevelopment bool  //Enable dev mode, all cache will be disabled
}

//New ivy with config
//source is the assets location
//cache is the location where cache will be store, or set nil if disable cache
func New(source Source, cache Cache, config *Config) *Ivy {
	return &Ivy{source, cache, config}
}

//Get get file or image
//if content type is image type and param not nil, process and return image
//if content type is image type and param nil, return raw image
//if content type is not image type, return raw file
//if cache enable, image type will cache. Next time return cache file
//if cache disable image type will not cahce, process and return image every time
//if HTTPCache is enable in config, return 304 status on next time
func (iv *Ivy) Get(bucket, paramsStr, path string, w http.ResponseWriter, r *http.Request) {
	if iv.isNotModify(r) {
		iv.writeNotModify(w, path)
		return
	}

	params, err := parseParams(paramsStr)
	if err != nil {
		iv.writeError(w, err)
		return
	}

	var img *bytes.Buffer

	if params.isDefault {
		if img, err = iv.loadFromSource(bucket, path, params); err != nil {
			iv.writeError(w, err)
			return
		}
	} else {
		if img, err = iv.loadFromCache(bucket, path, params); err != nil {
			if img, err = iv.loadFromSource(bucket, path, params); err != nil {
				iv.writeError(w, err)
				return
			}
		}
	}

	iv.writeSuccess(w, path, img)
}

//DeleteCache remove individual cache file
func (iv *Ivy) DeleteCache(bucket, filename string) error {
	return iv.Cache.Delete(bucket, filename)
}

//FlushCache remove all cache file in specific bucket
func (iv *Ivy) FlushCache(bucket string) error {
	return iv.Cache.Flush(bucket)
}

func (iv *Ivy) isNotModify(r *http.Request) bool {
	if iv.Config.HTTPCache > 0 && !iv.Config.IsDevelopment && r.Header.Get("If-Modified-Since") != "" {
		return true
	}

	return false
}

func (iv *Ivy) loadFromCache(bucket, filePath string, params *params) (*bytes.Buffer, error) {
	if iv.Config.IsDevelopment || iv.Cache == nil {
		return nil, errors.New("no cache")
	}

	file, err := iv.Cache.Load(bucket, filePath, params.String())
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (iv *Ivy) loadFromSource(bucket, filePath string, params *params) (*bytes.Buffer, error) {
	file, err := iv.Source.Load(bucket, filePath)
	if err != nil {
		return nil, err
	}

	img, err := process(params, iv.Source.GetFilePath(bucket, filePath), file)
	if err != nil {
		return nil, err
	}

	if iv.Cache != nil {
		go iv.Cache.Save(bucket, filePath, params.String(), img.Bytes())
	}

	return img, nil
}

func (iv *Ivy) getContentType(filePath string) string {
	switch filepath.Ext(filePath) {
	case ".jpg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	}

	return "application/octet-stream"
}

func (iv *Ivy) setHeader(w http.ResponseWriter, filePath string) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", iv.getContentType(filePath))
	w.Header().Add("Connection", "keep-alive")
	w.Header().Add("Vary", "Accept-Encoding")
	w.Header().Add("Last-Modified", "Tue, 01 Jan 2008 00:00:00 GMT")

	if iv.Config.HTTPCache > 0 && !iv.Config.IsDevelopment {
		w.Header().Add("Cache-Control", "public; max-age="+strconv.FormatInt(iv.Config.HTTPCache, 10))
		w.Header().Add("Expires", time.Now().Add(time.Second*time.Duration(iv.Config.HTTPCache)).Format("Mon, _2 Jan 2006 15:04:05 MST"))
	}
}

func (iv *Ivy) writeSuccess(w http.ResponseWriter, filePath string, file *bytes.Buffer) {
	iv.setHeader(w, filePath)

	w.WriteHeader(http.StatusOK)
	w.Write(file.Bytes())
}

func (iv *Ivy) writeError(w http.ResponseWriter, err error) {
	switch err {
	case ErrNotFound:
		w.WriteHeader(404)
		w.Write([]byte("not found"))
	default:
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}

func (iv *Ivy) writeNotModify(w http.ResponseWriter, filePath string) {
	iv.setHeader(w, filePath)
	w.WriteHeader(304)
	w.Write(nil)
}
