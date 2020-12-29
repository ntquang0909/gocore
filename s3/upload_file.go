package s3

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gabriel-vasile/mimetype"
)

// UploadFileParams upload params
type UploadFileParams struct {
	Data               []byte
	Bucket             string
	Key                string
	Metadata           map[string]*string
	ContentDisposition string
}

// UploadFile upload file
func (s3 *S3) UploadFile(params UploadFileParams) (string, error) {

	s, err := s3.NewSession()
	if err != nil {
		s3.logger.Printf("Create session err: %v\n", err)
		return "", err
	}

	var svc = s3manager.NewUploader(s)

	var mine = mimetype.Detect(params.Data)

	var uploadParams = &s3manager.UploadInput{
		Bucket:      aws.String(params.Bucket),
		Key:         aws.String(params.Key),
		ACL:         aws.String("public-read"),
		Body:        bytes.NewReader(params.Data),
		ContentType: aws.String(mine.String()),
	}
	if params.ContentDisposition != "" {
		uploadParams.ContentDisposition = aws.String(params.ContentDisposition)
	}
	if params.Metadata != nil {
		uploadParams.Metadata = params.Metadata
	}
	result, err := svc.Upload(uploadParams)
	if err != nil {
		s3.logger.Printf("Upload s3 err: %v\n", err)
		return "", err
	}

	return result.Location, nil

}
