package main

import (
	"go-initializr/app/initializer"
	"go-initializr/app/pkg"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
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
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Every(time.Minute / 10),
				Burst:     10,
				ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			log.Print("too many request ")
			return context.JSON(http.StatusTooManyRequests, map[string]any{
				"message": "too many request",
			})
		},
	}

	e.Use(middleware.RateLimiterWithConfig(config))
	e.Validator = pkg.NewCustomValidator()
	e.Use(middleware.Logger())
	e.GET("/hello-world", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Static("/assets", "frontend/assets") // Serve JS & CSS
	e.Static("/", "frontend")              // static assets
	e.GET("/*", func(c echo.Context) error {
		path := "frontend" + c.Request().URL.Path

		// Try to serve the file if it exists
		if _, err := os.Stat(path); err == nil {
			return c.File(path)
		}

		// Otherwise fallback to index.html for SPA
		return c.File("frontend/index.html")
	})

	version1 := e.Group("api/v1")
	handler := initializer.NewHandler()
	version1.POST("/initialize", handler.InitializeBoilerplate)
	version1.GET("/initialize/download/:id", handler.DownloadFolder)
	version1.POST("/reset-folder", handler.DeleteAllGeneratedProject)
	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))
}
