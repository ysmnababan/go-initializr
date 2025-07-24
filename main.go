package main

import (
	"errors"
	"go-initializr/app/initializer"
	"go-initializr/app/pkg"
	"go-initializr/app/pkg/logger"
	"go-initializr/app/pkg/response"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	zlog "github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

var APP_ENV string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, skipping...")
	}
	APP_ENV = "DEVELOPMENT"
	env := os.Getenv("ENV")
	if len(env) > 0 {
		APP_ENV = env
	}
	log.Println("ENV:", APP_ENV)
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
	if APP_ENV == "DEVELOPMENT" {
		e.Use(middleware.Logger())
	}
	logger.InitLogger()
	e.Use(logger.WithRequestLogger())
	e.HTTPErrorHandler = customHttpHandler
	e.GET("/hello-world", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Static("/assets", "frontend/dist/assets") // Serve JS & CSS
	e.Static("/", "frontend/dist")              // static assets
	e.File("/go-initializr-icon.svg", "frontend/dist/go-initializr-icon.svg")
	e.GET("/*", func(c echo.Context) error {
		return c.File("frontend/dist/index.html")
	})

	version1 := e.Group("api/v1")
	handler := initializer.NewHandler()
	version1.POST("/initialize", handler.InitializeBoilerplate)
	version1.GET("/initialize/download/:id", handler.DownloadFolder)
	version1.POST("/reset-folder", handler.DeleteAllGeneratedProject)
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323" // default for local/dev
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func customHttpHandler(err error, c echo.Context) {
	ctx := c.Request().Context()
	logger := zlog.Ctx(ctx)

	var apiErr *response.APIError
	if errors.As(err, &apiErr) {
		logger.Error().
			Err(err).
			Int("status_code", apiErr.Code).
			Msg(apiErr.Message)

		_ = c.JSON(apiErr.Code, response.APIResponse{
			Meta: response.Meta{
				Success:    false,
				Message:    apiErr.Message,
				StatusCode: apiErr.StatusCode,
				Detail:     apiErr.Detail,
			},
		})
		return
	}

	logger.Error().
		Err(err).
		Str("path", c.Path()).
		Msg("unhandled internal error")
	_ = c.JSON(http.StatusInternalServerError,
		response.APIResponse{
			Meta: response.Meta{
				Success:    false,
				Message:    response.ErrInternalServerError.Message,
				StatusCode: response.ErrInternalServerError.StatusCode,
			},
		})
}
