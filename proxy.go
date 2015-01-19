package fileproxy

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ErrNotFound = errors.New("not found")

type Source interface {
	Load(filename string) (io.Reader, int64, time.Time, error)
	GetFilePath(filename string) string
}

type Cache interface {
	Save(filename, paramsStr string, file []byte) error
	Load(filename, paramsStr string) (io.Reader, int64, time.Time, error)
	Delete(filename string) error
	Flush() error
}

type FileProxy struct {
	Source Source
	Cache  Cache
	Config *Config
	root   string
}

type Config struct {
	HttpCache     int64
	IsDevelopment bool
}

func New(root string, source Source, cache Cache, config *Config) *FileProxy {
	return &FileProxy{source, cache, config, root}
}

func (f *FileProxy) Load(w http.ResponseWriter, r *http.Request) {
	if f.Config.HttpCache > 0 && !f.Config.IsDevelopment && r.Header.Get("If-Modified-Since") != "" {
		f.writeNotModify(w)
		return
	}

	var err error
	var file io.Reader
	var size int64
	var modTime time.Time
	var filename string
	var paramsStr string
	path := strings.Replace(r.URL.Path, f.root, "", 1)
	strList := strings.Split(path, "/")

	switch len(strList) {
	case 2:
		filename = strList[1]
		paramsStr = ""
	case 3:
		filename = strList[2]
		paramsStr = strList[1]
	default:
		f.writeError(w, errors.New("Invalid Url"))
		return
	}

	params, err := parseParams(paramsStr)
	paramsStr = params.String()

	//load from cache
	if !f.Config.IsDevelopment {
		file, size, modTime, err = f.Cache.Load(filename, paramsStr)
		if err == nil {
			buffer := bytes.NewBuffer(nil)
			buffer.ReadFrom(file)
			f.writeSuccess(w, buffer.Bytes(), size, modTime)
			return
		}
	}

	//load from source
	if err == ErrNotFound {
		file, size, modTime, err = f.Source.Load(filename)
		if err == nil {
			img, err := process(params, f.Source.GetFilePath(filename), file)
			if err != nil {
				f.writeError(w, err)
				return
			}
			f.writeSuccess(w, img, size, modTime)
			if f.Cache != nil {
				f.Cache.Save(filename, paramsStr, img)
			}
			return
		}

		if err == ErrNotFound {
			f.writeNotFoud(w)
		} else {
			f.writeError(w, err)
		}
	}
}

func (f *FileProxy) FlushCache() error {
	return f.Cache.Flush()
}

func (f *FileProxy) writeSuccess(w http.ResponseWriter, file []byte, size int64, modTime time.Time) {
	w.Header().Add("Content-Type", http.DetectContentType(file))
	w.Header().Add("Connection", "keep-alive")
	w.Header().Add("Vary", "Accept-Encoding")
	w.Header().Add("Content-Length", strconv.FormatInt(size, 10))

	if f.Config.HttpCache > 0 && !f.Config.IsDevelopment {
		w.Header().Add("Cache-Control", "public; max-age="+strconv.FormatInt(f.Config.HttpCache, 10))
		w.Header().Add("Last-Modified", modTime.Format("Mon, _2 Jan 2006 15:04:05 MST"))
		w.Header().Add("Expires", time.Now().Add(time.Second*time.Duration(f.Config.HttpCache)).Format("Mon, _2 Jan 2006 15:04:05 MST"))
	}

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

func (f *FileProxy) writeNotModify(w http.ResponseWriter) {
	if f.Config.HttpCache > 0 && !f.Config.IsDevelopment {
		w.Header().Add("Cache-Control", "public; max-age="+strconv.FormatInt(f.Config.HttpCache, 10))
		w.Header().Add("Expires", time.Now().Add(time.Second*time.Duration(f.Config.HttpCache)).Format("Mon, _2 Jan 2006 15:04:05 MST"))
	}
	w.WriteHeader(304)
	w.Write(nil)
}
