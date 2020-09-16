package s3

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// UploadMultipleLogFileParams ignore file
type UploadMultipleLogFileParams struct {
	IgnoreFiles         []string
	FolderToUpload      string
	UploadToBucket      string
	KeepFileAfterUpload bool
}

// UploadMultipleLogFile update multiple log files
func (s3 *S3) UploadMultipleLogFile(params UploadMultipleLogFileParams) ([]string, error) {
	var response = []string{}
	var files []string
	var dir = params.FolderToUpload

	var ignoreFiles = params.IgnoreFiles

	var err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if len(ignoreFiles) > 0 {
			for _, ignoreFile := range ignoreFiles {
				if info.IsDir() == false && strings.Index(path, ignoreFile) == -1 {
					files = append(files, path)
				}
			}
		} else {
			if info.IsDir() == false {
				files = append(files, path)
			}
		}
		return nil
	})

	if err != nil {
		return response, err
	}

	if len(files) == 0 {
		log.Printf("No have any files to upload\n")
		return response, nil
	}

	var wg = sync.WaitGroup{}
	var max = len(files)
	var channel = make(chan string, max)
	for _, file := range files {
		wg.Add(1)
		go func(channel chan string, file string) error {
			defer func() {
				wg.Done()
				max--
				if max == 0 {
					close(channel)
				}

			}()

			log.Printf("Uploading file %s\n", file)

			originFile, err := os.Open(file)
			if err != nil {
				log.Printf("Open file: error %+v\n", err)
				return err
			}

			reader, writer := io.Pipe()
			go func() {
				gw := gzip.NewWriter(writer)
				io.Copy(gw, originFile)
				originFile.Close()
				gw.Close()
				writer.Close()
			}()
			var ext = path.Ext(file)
			var fileName = file[0 : len(file)-len(ext)]
			var gzipFileName = fmt.Sprintf("%s.gz", fileName)

			sess, err := s3.NewSession()
			if err != nil {
				log.Printf("Create s3 session error %v\n", err)
				return err
			}

			var fileKey = filepath.Base(gzipFileName)
			var folder = filepath.Dir(file)
			uploader := s3manager.NewUploader(sess)
			result, err := uploader.Upload(&s3manager.UploadInput{
				Body:   reader,
				Bucket: aws.String(params.UploadToBucket),
				Key:    aws.String(fmt.Sprintf("%s/%s", folder, fileKey)),
				ACL:    aws.String("public-read"),
			})
			if err != nil {
				log.Printf("Upload s3 error: %v\n", err)
				return err
			}

			log.Printf("%s is uploaded to s3 at %s\n", fileKey, result.Location)

			if params.KeepFileAfterUpload == false {
				err = os.Remove(file)
				if err != nil {
					log.Printf("Removed log file %s error %+v\n", file, err)
					return nil
				}
				log.Printf("Removed log file %s\n", file)
			}

			channel <- result.Location

			return nil
		}(channel, file)
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
