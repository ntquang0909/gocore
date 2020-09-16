package main

import "github.com/thaitanloi365/gocore/logger"

func printlog() {
	logger.Global().Debug("Call from global log")
}
