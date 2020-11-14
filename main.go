package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kjk/dailyrotate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/xid"
	"github.com/subosito/gotenv"
	"github.com/thaitanloi365/gocore/logger"
	"github.com/thaitanloi365/gocore/logger/notifier"
	"github.com/thaitanloi365/gocore/s3"
	"github.com/thaitanloi365/gocore/storage"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	gotenv.Load("./.env")
	testStorageServe()
	// testLogWithHTTP()
	// testS3()
}

func testLoggerWithDailyRotate() {
	writer, err := dailyrotate.NewFile("logs/2006-01-02.log", func(path string, didRotate bool) {})
	if err != nil {
		panic(err)
	}

	var logger = logger.New(&logger.Config{
		BufferedSize: 100,
		Writer:       log.New(writer, "", 0),
	})
	var data = []interface{}{
		"asdf", "ss", "sss",
	}

	logger.Debugf("%s\n[info] "+"asdf", append([]interface{}{"aaaaa"}, data...)...)
	for i := 0; i < 10; i++ {
		logger.Debugf("count %d \n", i)
		logger.Debug("count sssss", i, "asdfasdf")
		logger.DebugJSON("asdfasdfadsf", data, data)
		time.Sleep(time.Second)
	}
}

func testLoggerWithLumberjack() {
	var writer = &lumberjack.Logger{
		Filename:   "foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	var logger = logger.New(&logger.Config{
		BufferedSize: 100,
		Writer:       log.New(writer, "\r\n", 0),
	})
	var data = []interface{}{
		"asdf", "ss", "sss",
	}

	logger.Debugf("%s\n[info] "+"asdf", append([]interface{}{"aaaaa"}, data...)...)
	for i := 0; i < 10; i++ {
		logger.Debugf("count %d \n", i)
		logger.Debug("count sssss", i, "asdfasdf")
		time.Sleep(time.Second)
	}
}

func testStorageServe() {
	var e = echo.New()
	var st = storage.New(storage.DefaultConfig)

	st.NewRouter(e.Group("/api"))

	e.Start(":1234")
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
		fmt.Println(files)
		for _, file := range files {
			fmt.Println("file", file.Filename)
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

func testLogWithHTTP() {

	writer, err := dailyrotate.NewFile("logs/2006-01-02.log", func(path string, didRotate bool) {})
	if err != nil {
		panic(err)
	}

	var log = logger.New(&logger.Config{
		BufferedSize: 100,
		Writer:       log.New(writer, "", 0),
		Notifier: &notifier.SlackNotifier{
			WebhookURL: "",
			Channel:    "",
		},
	})

	var e = echo.New()
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: xid.New().String,
	}))
	e.GET("/success", func(c echo.Context) error {
		fmt.Println(c.Request().Header)
		c.Set(logger.RefErrorIDKey, "aaaaa")
		c.Set(logger.UserIDKey, "1111")
		var user = map[string]interface{}{
			"name": "Test",
		}
		log.DebugJSON(user)
		log.DebugJSONWithEchoContext(c, user)
		log.DebugfWithEchoContext(c, "Dispatch task id = %s name = %s success\n", "result.Signature.UUID", "result.Signature.Name")
		log.Debugf("Dispatch task id = %s name = %s success\n", "result.Signature.UUID", "result.Signature.Name")
		defer log.DebugJSONWithEchoContext(c, "Guess")
		return c.JSON(200, "Success")
	})

	e.GET("/error", func(c echo.Context) error {
		log.ErrorJSONWithEchoContext(c, "Guess")
		return echo.NewHTTPError(500, "test")
	})

	e.Start(":1234")
}
