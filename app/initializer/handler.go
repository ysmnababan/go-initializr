package initializer

import (
	"fmt"
	"go-initializr/app/pkg/response"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type InitializerService interface {
	InitializeBoilerplate(req *BasicConfigRequest) (folderId string, err error)
	DownloadProjectByFolderID(folderID string) (err error)
}

type handler struct {
	service InitializerService
}

func NewHandler() *handler {
	return &handler{
		service: NewService(),
	}
}

func (h *handler) InitializeBoilerplate(c echo.Context) (err error) {
	req := new(BasicConfigRequest)
	err = c.Bind(req)
	if err != nil {
		return response.ErrorWrap(response.ErrUnprocessableEntity, err).Send(c)
	}

	err = c.Validate(req)
	if err != nil {
		return response.ErrorWrap(response.ErrBadRequest, err).Send(c)
	}

	folderID, err := h.service.InitializeBoilerplate(req)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(
		map[string]any{
			"folder_id": folderID,
		}).Send(c)
}

func (h *handler) DownloadFolder(c echo.Context) (err error) {
	return nil
}

func (h *handler) DeleteAllGeneratedProject(c echo.Context) (err error) {
	entries, err := os.ReadDir(GENERATED_ROOT_FOLDER)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)

	}

	for _, entry := range entries {
		entryPath := filepath.Join(GENERATED_ROOT_FOLDER, entry.Name())
		err := os.RemoveAll(entryPath)
		if err != nil {
			return fmt.Errorf("failed to remove %s: %w", entryPath, err)
		}
	}
	return nil
}
