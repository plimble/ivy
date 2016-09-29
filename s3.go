package ivy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/plimble/errors"
)

type S3 struct {
	s3     *s3manager.Downloader
	bucket string
}

func NewS3(awsId, awsSecret, s3Bucket, s3Region string) Source {
	cred := credentials.NewStaticCredentials(awsId, awsSecret, "")
	session := session.New(&aws.Config{Region: aws.String(s3Region), Credentials: cred})

	s3manager.NewDownloader(session)

	return &S3{s3manager.NewDownloader(session), s3Bucket}
}

func (s *S3) Get(path string) ([]byte, error) {
	b := aws.NewWriteAtBuffer(nil)
	_, err := s.s3.Download(b, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		cerr := err.(awserr.Error)
		if cerr.Code() == "NoSuchKey" {
			return nil, errors.NotFound("not found")
		}

		return nil, errors.InternalServerErrorErr(err)
	}

	return b.Bytes(), nil
}
