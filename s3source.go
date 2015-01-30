package fileproxy

// import (
// 	"github.com/awslabs/aws-sdk-go/aws"
// 	"io"
// 	"os"
// 	"path"
// )

// type FileSystemSource struct {
// 	root string
// }

// func NewS3Source(root string) *FileSystemSource {
// 	return &FileSystemSource{root}
// }

// func (fs *FileSystemSource) Load(filename string) (io.Reader, error) {
// 	filename = path.Join(fs.root, filename)

// 	file, err := os.Open(filename)
// 	if os.IsNotExist(err) {
// 		return nil, ErrNotFound
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return file, nil
// }

// func (fs *FileSystemSource) GetFilePath(filename string) string {
// 	return path.Join(fs.root, filename)
// }
