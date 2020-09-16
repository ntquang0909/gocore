package s3

import (
	"io"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// UploadFileParams single file params
type UploadFileParams struct {
	Body               io.Reader
	Key                string
	ContentType        string
	ContentDisposition string
}

// UploadMultipleFileParams params
type UploadMultipleFileParams struct {
	UploadFiles    []UploadFileParams
	UploadToBucket string
	ACL            string
}

// UploadMultipleFileResponse response
type UploadMultipleFileResponse struct {
	Error    map[string]error
	Response map[string]string
}

// UploadMultipleFile upload multiple files
func (s3 *S3) UploadMultipleFile(params UploadMultipleFileParams) (*UploadMultipleFileResponse, error) {
	var response = &UploadMultipleFileResponse{
		Error:    make(map[string]error),
		Response: make(map[string]string),
	}

	session, err := s3.NewSession()
	if err != nil {
		return response, err
	}

	var wg = sync.WaitGroup{}
	var svc = s3manager.NewUploader(session)

	for _, param := range params.UploadFiles {
		wg.Add(1)
		go func(param *UploadFileParams, wg *sync.WaitGroup) {
			defer wg.Done()

			var acl = "public-read"
			var contentDisposition = "inline"
			if params.ACL != "" {
				acl = params.ACL
			}

			if param.ContentDisposition != "" {
				contentDisposition = param.ContentDisposition
			}

			var input = &s3manager.UploadInput{
				Bucket:             aws.String(params.UploadToBucket),
				Key:                aws.String(param.Key),
				Body:               param.Body,
				ACL:                aws.String(acl),
				ContentType:        aws.String(param.ContentType),
				ContentDisposition: aws.String(contentDisposition),
			}

			result, err := svc.Upload(input)
			if err != nil {

				s3.logger.Printf("Uploaded file %s at %s\n", param.Key, result.Location)
				response.Error[param.Key] = err
				return
			}

			response.Response[param.Key] = result.Location

		}(&param, &wg)
	}

	wg.Wait()

	return response, nil
}
