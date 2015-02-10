package ivy

import (
	"bytes"
)

type fakeSource struct {
	buffer *bytes.Buffer
	err    error
	path   string
}

func newFakeSource() *fakeSource {
	return &fakeSource{bytes.NewBuffer(nil), nil, ""}
}

func (s *fakeSource) Load(bucket string, filename string) (*bytes.Buffer, error) {
	return s.buffer, s.err
}

func (s *fakeSource) GetFilePath(bucket string, filename string) string {
	return s.path
}

type fakeCache struct {
	buffer *bytes.Buffer
	err    error
}

func newFakeCache() *fakeCache {
	return &fakeCache{bytes.NewBuffer(nil), nil}
}

func (c *fakeCache) Save(bucket, filename string, params *Params, file []byte) error {
	return c.err
}

func (c *fakeCache) Load(bucket, filename string, params *Params) (*bytes.Buffer, error) {
	return c.buffer, c.err
}

func (c *fakeCache) Delete(bucket, filename string, params *Params) error {
	return c.err
}

func (c *fakeCache) Flush(bucket string) error {
	return c.err
}

type fakeProcessor struct {
	buffer *bytes.Buffer
	err    error
}

func newFakeProcessor() *fakeProcessor {
	return &fakeProcessor{bytes.NewBuffer(nil), nil}
}

func (p *fakeProcessor) Process(params *Params, filePath string, file *bytes.Buffer) (*bytes.Buffer, error) {
	return p.buffer, p.err
}
