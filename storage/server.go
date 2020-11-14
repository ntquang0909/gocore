package storage

import (
	"io"
	"mime/multipart"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"github.com/thaitanloi365/gocore/cache"
)

// NewRouter setup router
func (storage *Storage) NewRouter(group *echo.Group) {
	group.POST("/upload", storage.uploadFileHandler)
	group.GET("/:id/download", storage.downloadFileHandler)
	group.DELETE("/:id/delete", storage.deleteFileHandler)
	group.GET("/list", storage.listFileHandler)
}

// UploadResult params
type UploadResult struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

func (storage *Storage) uploadFileHandler(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	var files = form.File["files"]
	var results []*UploadResult
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

			var id = xid.New().String()
			dst, err := storage.Create(id)
			if err != nil {
				storage.logger.Printf("Create file error: %v\n", err)
				return err
			}
			defer dst.Close()

			if _, err = io.Copy(dst, src); err != nil {
				storage.logger.Printf("Copy file error: %v\n", err)
				return err
			}

			if storage.config.Cache != nil {
				var cacheItem = FileCacheItem{
					ID:       id,
					Filename: file.Filename,
				}
				err = storage.config.Cache.Set(cacheItem.ID, &cacheItem, cache.NoExpiration)
				if err != nil {
					storage.logger.Printf("Set cache item error: %v\n", err)
				}
			}
			results = append(results, &UploadResult{
				ID:       id,
				Filename: file.Filename,
			})
			return nil
		}(file, &wg)
	}
	wg.Wait()

	var response = map[string]interface{}{
		"files": results,
	}
	return c.JSON(200, response)
}

func (storage *Storage) downloadFileHandler(c echo.Context) error {
	var fileID = c.Param("id")
	_, err := storage.Stat(fileID)
	if err != nil {
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
		return err
	}

	storage.Remove(fileID)

	var response = map[string]interface{}{
		"message": "Your file is deleted",
	}

	return c.JSON(200, response)
}
