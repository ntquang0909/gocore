package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/subosito/gotenv"
	"github.com/thaitanloi365/gocore/s3"
	"github.com/thaitanloi365/gocore/storage"
)

func main() {
	gotenv.Load("./.env")
	testS3()
}

func testS3() {
	var e = echo.New()

	e.PUT("/upload", func(c echo.Context) error {
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}

		var files = form.File["files"]
		var fileUploads = s3.UploadMultipleFileParams{
			ACL:            "public-read",
			UploadToBucket: os.Getenv("AWS_S3_ORIGIN_BUCKET"),
			UploadFiles:    []s3.UploadFileParams{},
		}
		for _, file := range files {
			fileUploads.UploadFiles = append(fileUploads.UploadFiles, s3.UploadFileParams{
				FileHeader: file,
				Prefix:     "test",
			})
		}
		var config = &s3.Config{
			AccessKey: os.Getenv("AWS_ACCESS_KEY"),
			SecretKey: os.Getenv("AWS_SECRET_KEY"),
			Region:    os.Getenv("AWS_REGION"),
		}
		result, err := s3.New(config).UploadMultipleFile(c.Request().Context(), fileUploads)
		if err != nil {
			return err
		}
		return c.JSON(200, result)
	})

	e.PUT("/upload_log", func(c echo.Context) error {
		var config = &s3.Config{
			AccessKey: os.Getenv("AWS_ACCESS_KEY"),
			SecretKey: os.Getenv("AWS_SECRET_KEY"),
			Region:    os.Getenv("AWS_REGION"),
		}
		result, err := s3.New(config).UploadMultipleLogFile(s3.UploadMultipleLogFileParams{
			UploadToBucket:      os.Getenv("AWS_S3_ORIGIN_BUCKET"),
			KeepFileAfterUpload: false,
			FolderToUpload:      "logs",
		})
		if err != nil {
			return err
		}
		return c.JSON(200, result)
	})
	e.Start(":1234")
}
func testStorage() {
	var storage = storage.New(&storage.Config{
		RootDir: "temp",
	})

	var e = echo.New()

	var fileGroup = e.Group("")
	fileGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte("secret"),
		TokenLookup: "query:token",
	}))
	fileGroup.Static("/file", "temp")

	go func() {
		for i := 0; i < 10; i++ {

			if i > 5 {
				storage.Create(fmt.Sprintf("/images/%d.png", i))
			} else {
				storage.Create(fmt.Sprintf("file_%d.csv", i))
			}
		}
	}()

	e.Start(":1234")
}
