package initializer

import "github.com/labstack/echo/v4"

type InitializerService interface {
	InitializeBoilerplate() (folderId string, err error)
}

type handler struct {
	service InitializerService
}

func NewHandler() *handler {
	return &handler{
		service: NewService(),
	}
}

func (h *handler) InitializeBoilerplate(e echo.Context) (err error) {
	return
}
