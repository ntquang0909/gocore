package storage

import (
	"fmt"
	"os"
	"path"

	"github.com/parnurzeal/gorequest"
)

// Client http client
type Client struct {
	*gorequest.SuperAgent
	baseURL string
	storage *Storage
}

// NewClient init
func newClient(storage *Storage, baseURL string) *Client {
	return &Client{
		SuperAgent: gorequest.New(),
		baseURL:    baseURL,
		storage:    storage,
	}
}

// UploadFileResult result
type UploadFileResult struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	DownloadURL string `json:"download_url"`
}

// UploadFileParams params
type UploadFileParams struct {
	Path string `json:"file_path"`
	Name string `json:"name"`
}

// UploadFiles upload file multiple files
func (c *Client) UploadFiles(params ...UploadFileParams) ([]*UploadFileResult, error) {
	var url = fmt.Sprintf("%s/upload", c.baseURL)
	var req = c.Post(url).Type("multipart")

	for _, param := range params {
		f, err := os.Open(param.Path)
		if err != nil {
			c.storage.logger.Printf("Open file error: %v", err)
			continue
		}

		if param.Name != "" && path.Ext(param.Name) == "" {
			param.Name = fmt.Sprintf("%s%s", param.Name, path.Ext(f.Name()))
		}
		req.SendFile(f, param.Name, "files")
	}

	var uploadedFiles []string
	_, _, errs := req.EndStruct(&uploadedFiles)
	if errs != nil && len(errs) > 0 {
		return nil, errs[0]
	}

	var results []*UploadFileResult
	for _, f := range uploadedFiles {
		var url = fmt.Sprintf("%s/%s", c.baseURL, f)
		results = append(results, &UploadFileResult{
			Name:        f,
			URL:         url,
			DownloadURL: fmt.Sprintf("%s/download", url),
		})
	}
	return results, nil
}

// DeleteFiles upload file multiple files
func (c *Client) DeleteFiles(files ...string) (errorsFiles []string) {
	for _, file := range files {
		var url = fmt.Sprintf("%s/%s", c.baseURL, file)
		_, _, errs := c.Delete(url).End()
		if errs != nil || len(errs) > 0 {
			c.storage.logger.Printf("Delete file %s error: %v", file, errs)
			errorsFiles = append(errorsFiles, file)
			continue
		}
	}

	return errorsFiles
}
