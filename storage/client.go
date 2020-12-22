package storage

import (
	"crypto/tls"
	"fmt"
	"os"
	"path"

	"github.com/parnurzeal/gorequest"
)

// BasicAuth basic auth
type BasicAuth struct {
	UserName string
	Password string
}

// ClientAuthConfig auth  config
type ClientAuthConfig struct {
	BasicAuth   *BasicAuth
	BearerToken string
}

// Client http client
type Client struct {
	*gorequest.SuperAgent
	baseURL string
	storage *Storage
	auth    *ClientAuthConfig
}

// NewClient init
func newClient(storage *Storage, baseURL string) *Client {
	var agent = gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	return &Client{
		SuperAgent: agent,
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
	req = c.setReqAuth(req)

	if c.auth != nil {
		if c.auth.BasicAuth != nil {
			req.SetBasicAuth(c.auth.BasicAuth.UserName, c.auth.BasicAuth.Password)
		}
		if c.auth.BearerToken != "" {
			req.Set("Authorization", fmt.Sprintf("Bearer %s", c.auth.BearerToken))
		}
	}

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
func (c *Client) DeleteFiles(files ...string) error {
	var url = fmt.Sprintf("%s", c.baseURL)
	var form = DeleteMultiFilesParams{
		Files: files,
	}
	var req = c.Delete(url)
	req = c.setReqAuth(req)

	_, _, errs := req.Send(form).End()
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (c *Client) setReqAuth(req *gorequest.SuperAgent) *gorequest.SuperAgent {
	if c.auth != nil {
		if c.auth.BasicAuth != nil {
			req.SetBasicAuth(c.auth.BasicAuth.UserName, c.auth.BasicAuth.Password)
		}
		if c.auth.BearerToken != "" {
			req.Set("Authorization", fmt.Sprintf("Bearer %s", c.auth.BearerToken))
		}
	}

	return req
}
