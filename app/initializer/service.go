package initializer

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"regexp"

	"go-initializr/app/pkg"
	"go-initializr/app/pkg/response"
	"os"
	"strings"

	"log"

	"github.com/google/uuid"
)

var validProjectName = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,64}$`)

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) InitializeBoilerplate(req *BasicConfigRequest) (zipData []byte, err error) {
	if len(req.ProjectName) > MAX_NAME_LENGTH || len(req.ModuleName) > MAX_NAME_LENGTH {
		err = response.ErrorWrap(response.ErrInputLength, errors.New("length exceeded 64 characters"))
		return
	}
	if err := validateProjectName(req.ProjectName); err != nil {
		err = response.ErrorWrap(response.ErrInvalidCharacter, err)
		return []byte{}, err
	}
	req.ProjectName = sanitizeProjectName(req.ProjectName)

	moduleName := sanitizeModuleName(req.ModuleName)
	log.Println("project,module:", req.ProjectName, moduleName)
	err = validateModuleName(moduleName)
	if err != nil {
		err = response.ErrorWrap(response.ErrInvalidCharacter, err)
		return
	}
	req.ModuleName = moduleName
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
	folderId := uuid.NewString()
	targetPath := fmt.Sprintf("%s/%s", GENERATED_ROOT_FOLDER, folderId)
	rootNode.Name = req.ProjectName
	err = rootNode.GenerateFolder(targetPath, req)
	defer func() {
		log.Println("deleting unused folder")
		err := os.RemoveAll(targetPath)
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return
	}

	rootProjectPath := fmt.Sprintf("%s/%s", targetPath, req.ProjectName)

	if err = runCommand(rootProjectPath, "go", "fmt", "./..."); err != nil {
		log.Println("Error running go format:", err)
		return
	}

	zipped, err := pkg.ZipFolder(rootProjectPath)
	if err != nil {
		err = response.ErrorWrap(response.ErrZip, err)
		return
	}
	return zipped, nil
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

// SanitizeModuleName transforms the input to a safer format
func sanitizeModuleName(input string) string {
	name := strings.TrimSpace(input)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")

	// Remove invalid characters
	reg := regexp.MustCompile(`[^a-z0-9\-\/\.]+`)
	name = reg.ReplaceAllString(name, "")

	// Normalize multiple dashes/slashes
	name = regexp.MustCompile(`[-]+`).ReplaceAllString(name, "-")
	name = regexp.MustCompile(`[/]+`).ReplaceAllString(name, "/")

	// Trim leading/trailing
	name = strings.Trim(name, "-/")
	return name
}

// ValidateModuleName checks if the module name is valid for go mod init
func validateModuleName(name string) error {
	if name == "" {
		return errors.New("module name cannot be empty")
	}

	if strings.Contains(name, " ") {
		return errors.New("module name must not contain spaces")
	}

	// Must start with a letter or domain-like part
	if !regexp.MustCompile(`^[a-z0-9]`).MatchString(name) {
		return errors.New("module name must start with a lowercase letter or number")
	}

	// Only allow lowercase letters, numbers, slashes, dots, and dashes
	if !regexp.MustCompile(`^[a-z0-9\-\/\.]+$`).MatchString(name) {
		return errors.New("module name contains invalid characters")
	}

	return nil
}

func sanitizeProjectName(input string) string {
	name := strings.TrimSpace(input)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")

	// Remove dangerous characters
	reg := regexp.MustCompile(`[^a-z0-9\-_]+`)
	name = reg.ReplaceAllString(name, "")

	// Normalize multiple dashes/underscores
	name = regexp.MustCompile(`[-_]+`).ReplaceAllString(name, "-")

	// Trim leading/trailing dashes
	name = strings.Trim(name, "-_")

	return name
}

func validateProjectName(name string) error {
	if name == "" {
		return errors.New("project name cannot be empty")
	}
	if len(name) < 3 {
		return errors.New("project name must be at least 3 characters")
	}
	if len(name) > 50 {
		return errors.New("project name must be less than 50 characters")
	}
	if !validProjectName.MatchString(name) {
		return errors.New("project name can only contain letters, numbers, hyphens, and underscores")
	}
	if name == "." || name == ".." {
		return errors.New("project name cannot be '.' or '..'")
	}
	return nil
}
