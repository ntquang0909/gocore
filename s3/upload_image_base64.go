package s3

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// ImageContentType content-type
type ImageContentType string

// Defined
var (
	PNG  ImageContentType = "image/png"
	JPEG ImageContentType = "image/jpeg"
)

// UploadImageBase64Params params
type UploadImageBase64Params struct {
	Key    string `json:"key"`
	Base64 string `json:"base64"`
	Bucket string `json:"bucket"`
}

// UploadImageBase64 base64
func (s3 *S3) UploadImageBase64(params UploadImageBase64Params) (string, error) {
	session, err := s3.NewSession()
	if err != nil {
		return "", err
	}

	rawData, err := base64.StdEncoding.DecodeString(params.GetRawBase64())
	if err != nil {
		return "", err
	}

	var buffer = bytes.NewReader(rawData)
	var contentType = params.GetContentType()
	var imageOrigin image.Image

	switch contentType {
	case PNG:
		imageOrigin, err = png.Decode(buffer)
		if err != nil {
			return "", err
		}
		break
	case JPEG:
		imageOrigin, err = jpeg.Decode(buffer)
		if err != nil {
			return "", err
		}
	}

	var ext = contentType.GetExtionsion()
	var size = imageOrigin.Bounds().Max
	var uploader = s3manager.NewUploader(session)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(params.Bucket),
		Key:         aws.String(fmt.Sprintf("%s%s", params.Key, ext)),
		Body:        bytes.NewBuffer(rawData),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(string(contentType)),
		Metadata: map[string]*string{
			"width":  aws.String(fmt.Sprintf("%d", size.X)),
			"height": aws.String(fmt.Sprintf("%d", size.Y)),
		},
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}

// GetContentType get content type
func (form *UploadImageBase64Params) GetContentType() ImageContentType {
	var from = strings.Index(string(form.Base64), ",")
	var suffix = strings.TrimSuffix(form.Base64[5:from], ";base64")
	switch suffix {

	case "image/png":
		return PNG
	case "image/jpeg":
		return JPEG
	case "image/jpg":
		return JPEG
	}
	return PNG
}

// GetRawBase64 get content type
func (form *UploadImageBase64Params) GetRawBase64() string {
	var b64data = form.Base64[strings.IndexByte(form.Base64, ',')+1:]
	return b64data
}

// GetExtionsion ext
func (imageContentType ImageContentType) GetExtionsion() string {
	switch imageContentType {
	case JPEG:
		return ".jpeg"
	default:
		break
	}

	return ".png"
}
