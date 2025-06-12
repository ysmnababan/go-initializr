package initializer

import (
	"bufio"
	"fmt"

	"go-initializr/app/pkg/response"
	"os"
	"strings"

	"github.com/google/uuid"
)

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) InitializeBoilerplate(req *BasicConfigRequest) (folderId string, err error) {
	// load the folder structure template
	file, err := os.Open(FOLDER_STRUCTURE_PATH)
	if err != nil {
		err = response.ErrorWrap(response.ErrOpeningFile, err)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	rootNode := &Node{}
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), " ")
		rootNode.ParseLine(line)
	}
	if err = scanner.Err(); err != nil {
		err = response.ErrorWrap(response.ErrScanner, err)
		return
	}

	// create the boilerplate project
	folderId = uuid.NewString()
	rootName := fmt.Sprintf("%s/%s", GENERATED_ROOT_FOLDER, folderId)
	rootNode.Name = req.ProjectName
	err = rootNode.GenerateFolder(rootName, req)
	if err != nil {
		return
	}

	return folderId, nil
}

func (s *service) DownloadProjectByFolderID(folderID string) (err error) {
	return
}
