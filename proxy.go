package fileproxy

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var ErrNotFound = errors.New("not found")

type Source interface {
	Load(filename string) (io.Reader, error)
	GetFilePath(filename string) string
}

type Cache interface {
	Save(filename, paramsStr string, file []byte) error
	Load(filename, paramsStr string) (io.Reader, error)
	Delete(filename string) error
	Flush() error
}

type FileProxy struct {
	Source Source
	Cache  Cache
	Config *Config
}

type Config struct {
	HttpCache     int64
	IsDevelopment bool
}

func New(source Source, cache Cache, config *Config) *FileProxy {
	return &FileProxy{source, cache, config}
}

func (f *FileProxy) Get(path string, w http.ResponseWriter, r *http.Request) {
	if f.isNotModify(r) {
		f.writeNotModify(w, path)
		return
	}

	filePath, params, err := f.parseURLPath(path)
	if err != nil {
		f.writeError(w, errors.New("Invalid Parse Url"))
		return
	}

	if img, err := f.loadFromCache(filePath, params); err == nil {
		f.writeSuccess(w, filePath, img)
		return
	}

	if img, err := f.loadFromSource(filePath, params); err == nil {
		f.writeSuccess(w, filePath, img)
		return
	} else {
		if err == ErrNotFound {
			f.writeNotFoud(w)
		} else {
			f.writeError(w, err)
		}
		return
	}
}

func (f *FileProxy) isNotModify(r *http.Request) bool {
	if f.Config.HttpCache > 0 && !f.Config.IsDevelopment && r.Header.Get("If-Modified-Since") != "" {
		return true
	}

	return false
}

func (f *FileProxy) loadFromCache(filePath string, params *Params) ([]byte, error) {
	if f.Config.IsDevelopment && f.Cache == nil {
		return nil, errors.New("no cache")
	}

	file, err := f.Cache.Load(filePath, params.String())
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)
	buffer.ReadFrom(file)

	return buffer.Bytes(), nil
}

func (f *FileProxy) loadFromSource(filePath string, params *Params) ([]byte, error) {
	file, err := f.Source.Load(filePath)
	if err != nil {
		return nil, err
	}

	img, err := process(params, f.Source.GetFilePath(filePath), file)
	if err != nil {
		return nil, err
	}

	if f.Cache != nil {
		go f.Cache.Save(filePath, params.String(), img)
	}

	return img, nil
}

func (f *FileProxy) parseURLPath(path string) (string, *Params, error) {
	var params *Params
	var filePath string
	var err error

	strList := strings.Split(path, "/")

	switch len(strList) {
	case 2:
		filePath = strList[1]
		params, err = parseParams("")
	case 3:
		filePath = strList[2]
		params, err = parseParams(strList[1])
	default:
		filePath = ""
		err = errors.New("Invalid Path Url")
	}

	return filePath, params, err
}

func (f *FileProxy) FlushCache() error {
	return f.Cache.Flush()
}

func (f *FileProxy) getContentType(filePath string) string {
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

func (f *FileProxy) setHeader(w http.ResponseWriter, filePath string) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", f.getContentType(filePath))
	w.Header().Add("Connection", "keep-alive")
	w.Header().Add("Vary", "Accept-Encoding")
	w.Header().Add("Last-Modified", "Tue, 01 Jan 2008 00:00:00 GMT")

	if f.Config.HttpCache > 0 && !f.Config.IsDevelopment {
		w.Header().Add("Cache-Control", "public; max-age="+strconv.FormatInt(f.Config.HttpCache, 10))
		w.Header().Add("Expires", time.Now().Add(time.Second*time.Duration(f.Config.HttpCache)).Format("Mon, _2 Jan 2006 15:04:05 MST"))
	}
}

func (f *FileProxy) writeSuccess(w http.ResponseWriter, filePath string, file []byte) {
	f.setHeader(w, filePath)

	w.WriteHeader(http.StatusOK)
	w.Write(file)
}

func (f *FileProxy) writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}

func (f *FileProxy) writeNotFoud(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Write([]byte("not found"))
}

func (f *FileProxy) writeNotModify(w http.ResponseWriter, filePath string) {
	f.setHeader(w, filePath)
	w.WriteHeader(304)
	w.Write(nil)
}
