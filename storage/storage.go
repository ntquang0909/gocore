package storage

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

var storage *Storage

// Storage storage
type Storage struct {
	config  *Config
	rootDir string
	logger  Logger
}

// FileInfo info
type FileInfo struct {
	Path       string
	Size       int64
	ModifyTime time.Time
}

// New storage
func New(cnf *Config) *Storage {
	var config = cnf
	if cnf == nil {
		config = DefaultConfig
	}

	var storage = &Storage{
		config:  config,
		rootDir: config.RootDir,
		logger:  DefaultConfig.Logger,
	}

	if config.Logger != nil {
		storage.logger = config.Logger
	}

	return storage

}

// Create new file
func (storage *Storage) Create(name string) (*os.File, error) {
	var path = storage.Path(name)
	var dir = filepath.Dir(path)

	var err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	var flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(path, flag, 0666)
	if err != nil {
		return file, err
	}
	_, err = file.Seek(0, io.SeekEnd)
	return file, err
}

// Open or create file if not exits
func (storage *Storage) Open(name string) (*os.File, error) {
	var path = storage.Path(name)
	var dir = filepath.Dir(path)

	var err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	var flag = os.O_RDWR | os.O_CREATE
	file, err := os.OpenFile(path, flag, 0666)
	if err != nil {
		return file, err
	}
	_, err = file.Seek(0, io.SeekEnd)
	return file, err
}

// Stat file info
func (storage *Storage) Stat(name string) (os.FileInfo, error) {
	var path = storage.Path(name)
	return os.Stat(path)

}

// Remove file
func (storage *Storage) Remove(names ...string) {
	for _, name := range names {
		var path = storage.Path(name)
		var err = os.Remove(path)
		if err != nil {
			storage.logger.Printf("Remove file error: %v\n", err)
		}
	}
}

// RemoveAll file
func (storage *Storage) RemoveAll(path string) error {
	var fullpath = storage.Path(path)
	return os.RemoveAll(fullpath)

}

// Walk walk all file
func (storage *Storage) Walk() ([]*FileInfo, error) {
	var files = []*FileInfo{}
	var err = filepath.Walk(storage.rootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		var fileInfo = FileInfo{
			Path:       path,
			Size:       info.Size(),
			ModifyTime: info.ModTime(),
		}
		files = append(files, &fileInfo)

		return nil
	})

	if err != nil {
		storage.logger.Printf("Walk error: %v\n", err)
	}
	return files, err
}

// Path get full path with root dir
func (storage *Storage) Path(name string) string {
	return filepath.Join(storage.rootDir, name)
}
