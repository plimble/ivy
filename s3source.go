package ivy

import (
	"bytes"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/s3"
	"github.com/plimble/utils/errors2"
	"strings"
)

//S3Source is aws s3 asset source
type S3Source struct {
	cli *s3.S3
}

//NewS3Source create s3 source with access key and secret key
func NewS3Source(accessKey, secretKey string) *S3Source {
	cli := s3.New(&aws.Config{
		Credentials: aws.Creds(accessKey, secretKey, ""),
		Region:      "ap-southeast-1",
	})

	return &S3Source{cli}
}

//Load s3 file
func (fs *S3Source) Load(bucket, filename string) (*bytes.Buffer, error) {
	if strings.HasPrefix(filename, "/") {
		filename = filename[1:]
	}

	res, err := fs.cli.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &filename,
	})

	if err != nil {
		if err.Error() == "Access Denied" || err.Error() == "The specified key does not exist." {
			return nil, errors2.NewNotFound("not found")
		}

		return nil, err
	}

	defer res.Body.Close()

	buffer := bytes.NewBuffer(make([]byte, 0, *res.ContentLength))
	buffer.ReadFrom(res.Body)

	return buffer, nil
}

//GetFilePath return filepath
func (fs *S3Source) GetFilePath(bucket, filename string) string {
	return filename
}
