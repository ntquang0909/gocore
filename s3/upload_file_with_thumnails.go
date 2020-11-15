package s3

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"path"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
	"github.com/rs/xid"
)

// ThumbnailSize size
type ThumbnailSize struct {
	Width  int
	Height int
	Prefix string
}

// UploadFileWithThumbnailParams single file params
type UploadFileWithThumbnailParams struct {
	FileHeader         *multipart.FileHeader
	Prefix             string
	ContentDisposition string
	ThumbnailSize      *ThumbnailSize
}

// UploadMultipleFileWithThumbnailParams params
type UploadMultipleFileWithThumbnailParams struct {
	UploadFiles             []UploadFileWithThumbnailParams
	UploadToBucket          string
	UploadToThumbnailBucket string
	ACL                     string
}

type UploadFileWithThumbnailResponse struct {
	Photo     string
	Thumbnail string
}

// CalcResponsiveSize calc responsive size
func (size *ThumbnailSize) CalcResponsiveSize(width, height int) {
	if size.Width == 0 && size.Height == 0 {
		size.Width = 128
	}

	if size.Width <= 0 {
		size.Width = size.Height * width / height
	}

	if size.Height <= 0 {
		size.Height = size.Width * height / width
	}

}

// UploadMultipleFileWithThumbnail upload multiple files
func (s3 *S3) UploadMultipleFileWithThumbnail(ctx context.Context, params UploadMultipleFileWithThumbnailParams) ([]*UploadFileWithThumbnailResponse, error) {
	var response = []*UploadFileWithThumbnailResponse{}

	session, err := s3.NewSession()
	if err != nil {
		return response, err
	}

	var wg = sync.WaitGroup{}
	var svc = s3manager.NewUploader(session)
	var max = len(params.UploadFiles)
	var channel = make(chan *UploadFileWithThumbnailResponse, max)
	wg.Add(max)
	for _, p := range params.UploadFiles {
		go func(channel chan *UploadFileWithThumbnailResponse, param UploadFileWithThumbnailParams, wg *sync.WaitGroup) {
			defer func() {
				wg.Done()
				max--
				if max == 0 {
					close(channel)
				}

			}()

			file, err := param.FileHeader.Open()
			if err != nil {
				channel <- nil
				return
			}
			defer file.Close()

			var fileID = xid.New()
			var ct = param.FileHeader.Header.Get("Content-Type")
			var ext = path.Ext(param.FileHeader.Filename)
			var key = fmt.Sprintf("%s/%s%s", param.Prefix, fileID, ext)

			s3.logger.Printf("Uploading file %s\n", key)

			var acl = "public-read"
			var contentDisposition = "inline"
			if params.ACL != "" {
				acl = params.ACL
			}

			if param.ContentDisposition != "" {
				contentDisposition = param.ContentDisposition
			}

			var imageOrigin image.Image
			var contentType = ImageContentType(ct)
			switch contentType {
			case PNG:
				imageOrigin, err = png.Decode(file)
				if err != nil {
					channel <- nil
					return
				}

			case JPEG, JPG:
				imageOrigin, err = jpeg.Decode(file)
				if err != nil {
					channel <- nil
					return
				}
			}

			var size = imageOrigin.Bounds().Max
			result, err := svc.Upload(&s3manager.UploadInput{
				Bucket:             aws.String(params.UploadToBucket),
				Key:                aws.String(key),
				Body:               file,
				ACL:                aws.String(acl),
				ContentType:        aws.String(string(contentType)),
				ContentDisposition: aws.String(contentDisposition),
				Metadata: map[string]*string{
					"width":  aws.String(fmt.Sprintf("%d", size.X)),
					"height": aws.String(fmt.Sprintf("%d", size.Y)),
				},
			})
			if err != nil {
				channel <- nil
				return
			}

			var uploadResult = &UploadFileWithThumbnailResponse{
				Photo: result.Location,
			}
			s3.logger.Printf("Upload file %s at %s\n", key, result.Location)

			if param.ThumbnailSize != nil {
				param.ThumbnailSize.CalcResponsiveSize(size.X, size.Y)
				var thumbnail = imaging.Resize(imageOrigin, param.ThumbnailSize.Width, param.ThumbnailSize.Height, imaging.Lanczos)
				var format imaging.Format
				switch contentType {
				case JPEG, JPG:
					format = imaging.JPEG
				case PNG:
					format = imaging.PNG
				}

				var bufferEncode = new(bytes.Buffer)
				err = imaging.Encode(bufferEncode, thumbnail, format)
				if err == nil {
					var key = fmt.Sprintf("%s/%s%s", param.Prefix, fileID, ext)
					result, err := svc.Upload(&s3manager.UploadInput{
						Bucket:             aws.String(params.UploadToThumbnailBucket),
						Key:                aws.String(key),
						Body:               bufferEncode,
						ACL:                aws.String(acl),
						ContentType:        aws.String(string(contentType)),
						ContentDisposition: aws.String(contentDisposition),
						Metadata: map[string]*string{
							"width":  aws.String(fmt.Sprintf("%d", param.ThumbnailSize.Width)),
							"height": aws.String(fmt.Sprintf("%d", param.ThumbnailSize.Height)),
						},
					})
					if err == nil {
						uploadResult.Thumbnail = result.Location
						s3.logger.Printf("Upload thumbnail file %s at %s\n", key, result.Location)
					}
				}
			}

			channel <- uploadResult

		}(channel, p, &wg)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			value, more := <-channel
			if more {
				if value != nil {
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
