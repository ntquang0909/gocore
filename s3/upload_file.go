package s3

import (
	"context"
	"fmt"
	"mime/multipart"
	"path"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rs/xid"
)

// UploadFileParams single file params
type UploadFileParams struct {
	FileHeader         *multipart.FileHeader
	Prefix             string
	ContentDisposition string
}

// UploadMultipleFileParams params
type UploadMultipleFileParams struct {
	UploadFiles    []UploadFileParams
	UploadToBucket string
	ACL            string
}

// UploadMultipleFile upload multiple files
func (s3 *S3) UploadMultipleFile(ctx context.Context, params UploadMultipleFileParams) ([]string, error) {
	var response = []string{}

	session, err := s3.NewSession()
	if err != nil {
		return response, err
	}

	var wg = sync.WaitGroup{}
	var svc = s3manager.NewUploader(session)
	var max = len(params.UploadFiles)
	var channel = make(chan string, max)

	for _, param := range params.UploadFiles {
		wg.Add(1)
		go func(channel chan string, param *UploadFileParams, wg *sync.WaitGroup) {
			defer func() {
				wg.Done()
				max--
				if max == 0 {
					close(channel)
				}

			}()

			file, err := param.FileHeader.Open()
			if err != nil {
				channel <- ""
				return
			}
			defer file.Close()

			var contentType = param.FileHeader.Header.Get("Content-Type")
			var ext = path.Ext(param.FileHeader.Filename)
			var key = fmt.Sprintf("%s/%s%s", param.Prefix, xid.New(), ext)

			s3.logger.Printf("Uploading file %s\n", key)

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
				Key:                aws.String(key),
				Body:               file,
				ACL:                aws.String(acl),
				ContentType:        aws.String(contentType),
				ContentDisposition: aws.String(contentDisposition),
			}

			result, err := svc.Upload(input)
			if err != nil {
				s3.logger.Printf("Upload file %s error: %v\n", key, err)
				channel <- ""
				return
			}

			s3.logger.Printf("Upload file %s at %s\n", key, result.Location)

			channel <- result.Location

		}(channel, &param, &wg)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			value, more := <-channel
			if more {
				if value != "" {
					response = append(response, value)
				}
			} else {
				return
			}
		}
	}()

	wg.Wait()

	return response, nil
}
