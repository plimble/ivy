package ivy

import (
	"bytes"
	"os"
	"path"
)

type FileSystemSource struct {
	root string
}

func NewFileSystemSource(root string) *FileSystemSource {
	return &FileSystemSource{root}
}

func (fs *FileSystemSource) Load(bucket string, filename string) (*bytes.Buffer, error) {
	filename = path.Join(fs.root, bucket, filename)

	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := bytes.NewBuffer(nil)
	buffer.ReadFrom(file)

	return buffer, nil
}

func (fs *FileSystemSource) GetFilePath(bucket, filename string) string {
	return path.Join(fs.root, bucket, filename)
}
