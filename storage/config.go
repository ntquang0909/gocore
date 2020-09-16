package storage

import (
	"log"
	"os"
)

// DefaultConfig default config
var DefaultConfig = &Config{
	RootDir: "/tmp/gocore/storage",
	Logger:  log.New(os.Stdout, "\r\n", 0),
}

// Config  config
type Config struct {
	RootDir string
	Logger  Logger
}
