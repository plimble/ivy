package fileproxy

import (
	"bytes"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/gen/s3"
	"github.com/plimble/errs"
	"io"
	"strings"
)

type S3Source struct {
	cli *s3.S3
}

func NewS3Source(accessKey, secretKey string) *S3Source {
	creds := aws.Creds(accessKey, secretKey, "")
	cli := s3.New(creds, "ap-southeast-1", nil)

	return &S3Source{cli}
}

func (fs *S3Source) Load(bucket, filename string) (io.Reader, error) {
	if strings.HasPrefix(filename, "/") {
		filename = filename[1:]
	}

	res, err := fs.cli.GetObject(&s3.GetObjectRequest{
		Bucket: &bucket,
		Key:    &filename,
	})

	if err != nil {
		if err.Error() == "Access Denied" {
			return nil, errs.NewNotFound("not found")
		}

		return nil, err
	}

	buffer := bytes.NewBuffer(nil)
	buffer.ReadFrom(res.Body)

	return buffer, nil
}

func (fs *S3Source) GetFilePath(bucket, filename string) string {
	return filename
}
