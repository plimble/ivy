package ivy

import (
	"bytes"
	"os"
	"path"
)

//FileSystemSource is file system source
type FileSystemSource struct {
	root string
}

//NewFileSystemSource create new file source
func NewFileSystemSource(root string) *FileSystemSource {
	return &FileSystemSource{root}
}

//Load file from file system
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

//GetFilePath get file path
func (fs *FileSystemSource) GetFilePath(bucket, filename string) string {
	return path.Join(fs.root, bucket, filename)
}
