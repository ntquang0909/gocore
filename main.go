package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/thaitanloi365/gocore/storage"
)

func main() {

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
