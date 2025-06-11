package pkg

import "fmt"

type Node struct {
	Name     string
	Level    int
	IsFolder bool

	Content    []*Node
	ParentNode *Node
}

func (n *Node) ParseLine(line string) {
	// line can't be empty
	if len(line) == 0 {
		return
	}
	runes := []Rune(line)
	keyword := []Rune{}

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
	if currentNode.IsRoot() {
		parent.IsFolder = true
		parent.Content = nil
		parent.Level = 0
		parent.Name = currentNode.Name
		parent.ParentNode = nil
		return
	}
	if parent.IsParentOf(currentNode) {
		currentNode.ParentNode = parent
		parent.Content = append(parent.Content, currentNode)
		return
	}
	if parent.IsLevelHigherThan(currentNode) {
		lastChild := parent.Content[len(parent.Content)-1]
		lastChild.insertNode(currentNode)
		return
	}
}

func (n *Node) IsRoot() bool {
	return n.Level == 0
}
func (parent *Node) IsLevelHigherThan(n *Node) bool {
	return parent.Level < n.Level
}

func (parent *Node) IsParentOf(n *Node) bool {
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
