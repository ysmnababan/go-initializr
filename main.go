package main

import (
	"go-initializr/app/initializer"
	"go-initializr/app/pkg"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, skipping...")
	}
}

func main() {
	if err := os.MkdirAll(initializer.GENERATED_ROOT_FOLDER, os.ModePerm); err != nil {
		panic(err)
	}
	e := echo.New()
	e.Validator = pkg.NewCustomValidator()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	version1 := e.Group("/v1")
	handler := initializer.NewHandler()
	version1.POST("/initialize", handler.InitializeBoilerplate)
	version1.GET("/initialize/download/:id", handler.DownloadFolder)
	version1.POST("/reset-folder", handler.DeleteAllGeneratedProject)
	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))
}
