package cayl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
)

const (
	Source = "source"
	Cache  = "cache"

	File = "file"
	S3   = "s3"
)

var (
	NotFound = errors.New("not found")
)

func IsNotFound(err error) bool {
	if err != NotFound {
		return false
	}

	return true
}

type Storage interface {
	Save(filename string, data []byte) error
	Load(filename string) (io.Reader, error)
	Delete(filename string) error
	DeleteFolder(dir string) error
	Exist(path string) (bool, error)
	IsNotFound(err error) bool
}

type Config struct {
	FileStoragePath string
}

type Proxy struct {
	File Storage
	S3   Storage
}

func New(config Config) *Proxy {

	if config.FileStoragePath == "" {
		config.FileStoragePath = "./"
	}

	return &Proxy{
		File: newFileStorage(config.FileStoragePath),
	}
}

func (p *Proxy) getStorage(storage string) Storage {
	switch storage {
	case File:
		return p.File
	case S3:
		return p.S3
	default:
		return p.File
	}
}

func (p *Proxy) Save(storage, filename string, data multipart.File) (*bytes.Buffer, error) {
	s := p.getStorage(storage)

	buffer := bytes.NewBuffer(nil)
	buffer.ReadFrom(data)

	exist, err := s.Exist(sourcePath(filename))
	if exist {
		return nil, fmt.Errorf("%s is already exist", filename)
	}
	if err != nil {
		return nil, err
	}

	err = s.Save(sourcePath(filename), buffer.Bytes())
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func (p *Proxy) Load(storage, paramsStr, filename, ext string) ([]byte, error) {
	s := p.getStorage(storage)
	params, err := ParseParams(paramsStr)
	if err != nil {
		return nil, err
	}
	sp := sourcePath(filename)
	cp := cachePath(params, filename, ext)

	//load default
	if params.IsDefault {
		reader, err := s.Load(sp)
		if err != nil {
			if s.IsNotFound(err) {
				return nil, NotFound
			}
			return nil, err
		}
		buffer := bytes.NewBuffer(nil)
		buffer.ReadFrom(reader)

		return buffer.Bytes(), nil
	}

	//load from cache
	reader, err := s.Load(cp)
	if err != nil && !s.IsNotFound(err) {
		return nil, err
	}
	if err == nil {
		buffer := bytes.NewBuffer(nil)
		buffer.ReadFrom(reader)

		return buffer.Bytes(), nil
	}

	//load source for process
	reader, err = s.Load(sp)
	if err != nil {
		if s.IsNotFound(err) {
			return nil, NotFound
		}
		return nil, err
	}

	data, err := process(params, sp, ext, reader)
	if err != nil {
		return nil, err
	}

	go s.Save(cp, data)

	return data, nil
}

func (p *Proxy) Delete(storage, filename string) error {
	s := p.getStorage(storage)

	if err := s.Delete(sourcePath(filename)); err != nil {
		return err
	}

	return s.DeleteFolder(Cache + "/" + filename)
}

func cachePath(params *Params, filename, ext string) string {
	return Cache + "/" + filename + "/" + params.String() + ext
}

func sourcePath(filename string) string {
	return Source + "/" + filename
}
