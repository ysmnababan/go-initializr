package initializer

import (
	"fmt"
	"go-initializr/app/pkg"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Node struct {
	Name     string
	Level    int
	IsFolder bool

	Content    []*Node
	ParentNode *Node
}

func (n *Node) GenerateFile(root string, config *BasicConfigRequest) (err error) {
	tmplFilename, ok := TEMPLATE_REGISTRY[n.Name]
	if !ok {
		log.Println("template registry map not found: ", n.Name)
		return
	}

	if strings.Contains(tmplFilename, "redis") && !config.Redis {
		return
	}
	if strings.Contains(tmplFilename, "swagger") && !config.Swagger {
		return
	}
	if strings.Contains(tmplFilename, "auth") && !config.JWT {
		return
	}
	if strings.Contains(tmplFilename, "validator") && !config.Validator {
		return
	}
	if strings.Contains(tmplFilename, "token") && !config.JWT {
		return
	}
	tmplPath := filepath.Join(TEMPLATE_FOLDER_PATH, tmplFilename) // Assuming templates are in a "templates" directory
	tmpl := template.Must(template.New(tmplFilename).ParseFiles(tmplPath))
	filename := filepath.Join(root, n.Name)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error when creating file: %w", err)
	}

	err = tmpl.Execute(file, config)
	if err != nil {
		return fmt.Errorf("error when executing template: %w", err)
	}
	return nil
}

func (n *Node) GenerateFolder(root string, config *BasicConfigRequest) error {
	if n.IsFolder {
		currentPath := filepath.Join(root, n.Name)
		err := os.MkdirAll(currentPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error when creating folder: %w", err)
		}
		if len(n.Content) > 0 {
			for _, child := range n.Content {
				err := child.GenerateFolder(currentPath, config)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	// create file
	err := n.GenerateFile(root, config)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) ParseLine(line string) {
	// line can't be empty
	if len(line) == 0 {
		return
	}
	runes := []pkg.Rune(line)
	keyword := []pkg.Rune{}

	node := new(Node)
	if runes[len(runes)-1].IsSlash() {
		node.IsFolder = true
	}
	i := 0
	// for space parsing
	for {
		if !runes[i].IsSpace() || i >= len(runes)-1 {
			break
		}
		i++

		node.Level++
	}

	// store the keyword
	for {
		if i == len(runes)-1 && runes[i].IsSlash() {
			break
		}

		keyword = append(keyword, runes[i])
		i++
		if i >= len(runes) {
			break
		}
	}
	node.Level = node.Level / 2
	node.Name = string(keyword)
	n.insertNode(node)
}

func (parent *Node) insertNode(currentNode *Node) {
	if currentNode.isRoot() {
		parent.IsFolder = currentNode.IsFolder
		parent.Content = nil
		parent.Level = 0
		parent.Name = currentNode.Name
		parent.ParentNode = nil
		return
	}
	if parent.isParentOf(currentNode) {
		currentNode.ParentNode = parent
		parent.Content = append(parent.Content, currentNode)
		return
	}
	if parent.isLevelHigherThan(currentNode) {
		lastChild := parent.Content[len(parent.Content)-1]
		lastChild.insertNode(currentNode)
		return
	}
}

func (n *Node) isRoot() bool {
	return n.Level == 0
}
func (parent *Node) isLevelHigherThan(n *Node) bool {
	return parent.Level < n.Level
}

func (parent *Node) isParentOf(n *Node) bool {
	return parent.Level+1 == n.Level
}

func (n *Node) PrintNode() {
	eol := ""
	if n.IsFolder {
		eol = "/"
	}
	for range n.Level * 2 {
		fmt.Print("_")
	}

	fmt.Println(n.Name + eol)
	if !n.IsFolder || len(n.Content) == 0 {
		return
	}

	for _, child := range n.Content {
		child.PrintNode()
	}
}
