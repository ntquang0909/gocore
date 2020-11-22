package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
)

// WithRouter with router
func (storage *Storage) WithRouter(group *echo.Group) *Storage {
	group.POST("/upload", storage.uploadFileHandler)
	group.GET("/:id/download", storage.downloadFileHandler)
	group.GET("/:id", storage.previewFileHandler)
	group.DELETE("/:id", storage.deleteFileHandler)
	group.DELETE("", storage.deleteMultiFilesHandler)
	group.GET("/list", storage.listFileHandler)
	group.Static("/images", storage.rootDir)

	return storage
}

// WithClient with client
func (storage *Storage) WithClient(baseURL string, authConfig *ClientAuthConfig) *Storage {
	storage.client = newClient(storage, baseURL)
	storage.client.auth = authConfig
	return storage
}

func (storage *Storage) uploadFileHandler(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	var files = form.File["files"]
	var results []string
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader, wg *sync.WaitGroup) error {
			defer wg.Done()

			src, err := file.Open()
			if err != nil {
				storage.logger.Printf("Open file error: %v\n", err)
				return err
			}
			defer src.Close()

			storage.logger.Printf("Uploading file %s size = %v\n", file.Filename, file.Size)

			dst, err := storage.Create(file.Filename)
			if err != nil {
				storage.logger.Printf("Create file error: %v\n", err)
				return err
			}
			defer dst.Close()

			if _, err = io.Copy(dst, src); err != nil {
				storage.logger.Printf("Copy file error: %v\n", err)
				return err
			}

			results = append(results, file.Filename)
			return nil
		}(file, &wg)
	}
	wg.Wait()

	return c.JSON(200, results)
}

func (storage *Storage) downloadFileHandler(c echo.Context) error {
	var fileID = c.Param("id")
	_, err := storage.Stat(fileID)
	if err != nil {
		if os.IsNotExist(err) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not found", fileID))
		}

		return err
	}

	var file = storage.Path(fileID)
	var name = fileID

	if storage.config.Cache != nil {
		var cacheItem = FileCacheItem{}
		err = storage.config.Cache.Get(fileID, &cacheItem)
		if err != nil {
			storage.logger.Printf("Set cache item error: %v\n", err)
		} else {
			name = cacheItem.Filename
		}
	}
	return c.Attachment(file, name)
}

func (storage *Storage) previewFileHandler(c echo.Context) error {
	var fileID = c.Param("id")
	_, err := storage.Stat(fileID)
	if err != nil {
		if os.IsNotExist(err) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not found", fileID))
		}

		return err
	}

	var file = storage.Path(fileID)
	// var name = fileID

	return c.File(file)
}

func (storage *Storage) listFileHandler(c echo.Context) error {
	files, err := storage.Walk()
	if err != nil {
		return err
	}
	var response = map[string]interface{}{
		"files": files,
	}
	return c.JSON(200, response)
}

func (storage *Storage) deleteFileHandler(c echo.Context) error {
	var fileID = c.Param("id")
	_, err := storage.Stat(fileID)
	if err != nil {
		if os.IsNotExist(err) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not found", fileID))
		}

		storage.logger.Printf("Delete file error: %v\n", err)
		return err
	}

	storage.logger.Printf("Deleting file %v\n", fileID)
	storage.Remove(fileID)

	var response = map[string]interface{}{
		"message": "Your file is deleted",
	}

	return c.JSON(200, response)
}

// DeleteMultiFilesParams params
type DeleteMultiFilesParams struct {
	Files []string `json:"files"`
}

func (storage *Storage) deleteMultiFilesHandler(c echo.Context) error {
	var form DeleteMultiFilesParams
	var err = c.Bind(&form)
	if err != nil {
		return err
	}

	storage.Remove(form.Files...)

	return c.JSON(200, nil)
}
