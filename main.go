package main

import (
	"go-initializr/app/initializer"
	"go-initializr/app/pkg"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Validator = pkg.NewCustomValidator()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	version1 := e.Group("/v1")
	handler := initializer.NewHandler()
	version1.POST("/initialize", handler.InitializeBoilerplate)
	version1.GET("/initialize/download/:id", handler.DownloadFolder)
	e.Logger.Fatal(e.Start(":1323"))
}
