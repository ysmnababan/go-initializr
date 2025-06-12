package initializer

import (
	"bufio"
	"fmt"
	"os/exec"

	"go-initializr/app/pkg/response"
	"os"
	"strings"

	"log"

	"github.com/google/uuid"
)

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) InitializeBoilerplate(req *BasicConfigRequest) (folderId string, err error) {
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

	projectDir := fmt.Sprintf("%s/%s", rootName, req.ProjectName)
	fmt.Println("Initializing Go module in", projectDir)
	if err = runCommand(projectDir, "go", "mod", "init", req.ProjectName); err != nil {
		log.Println("Error initializing module:", err)
		return
	}

	if err = runCommand(projectDir, "go", "fmt", "./..."); err != nil {
		log.Println("Error running go format:", err)
		return
	}
	return folderId, nil
}

func runCommand(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir // Set working directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *service) DownloadProjectByFolderID(folderID string) (err error) {
	return
}
