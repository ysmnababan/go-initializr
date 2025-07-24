package initializer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLine_RealisticTree(t *testing.T) {
	t.Parallel()
	root := &Node{}
	lines := []string{
		"root/",
		"  cmd/",
		"    cron/",
		"  config/",
		"    config.go",
		"    stringutils.go",
		"    config_models.go",
		"  internal/",
		"    abstraction/",
		"      entity.go",
		"      pagination.go",
		"    factory/",
		"      factory.go",
		"    app/",
		"      example_feat/",
		"        controller.go",
		"        dto.go",
		"        model.go",
		"        repository.go",
		"        routes.go",
		"        service.go",
		"      next_feat/",
		"    middleware/",
		"      auth.go",
		"      middleware.go",
		"    server/",
		"      server_routes.go",
		"    pkg/",
		"      database/",
		"        database.go",
		"      redisutil/",
		"        redis.go",
		"    utils/",
		"      env/",
		"        env.go",
		"      response/",
		"        error_handler.go",
		"        response_model.go",
		"        success_handler.go",
		"      token/",
		"        jwt.go",
		"      validator/",
		"        custom_validator.go",
		"        validator.go",
		"  demo/",
		"  migrations/",
		"  test/",
		"  Dockerfile",
		"  .env",
		"  .gitignore",
		"  README.md",
		"  main.go",
		"  go.mod",
	}

	for _, line := range lines {
		root.ParseLine(line)
	}

	// Now assert selected paths exist:
	assert.Equal(t, "root", root.Name)
	assert.True(t, root.IsFolder)

	// Internals/app/example_feat/controller.go
	exampleFeat := findNodeByPath(root, "internal", "app", "example_feat")
	require.NotNil(t, exampleFeat)
	assert.Equal(t, "example_feat", exampleFeat.Name)
	assert.True(t, exampleFeat.IsFolder)

	controller := findNodeByPath(exampleFeat, "controller.go")
	require.NotNil(t, controller)
	assert.Equal(t, "controller.go", controller.Name)
	assert.False(t, controller.IsFolder)

	// Root-level file
	mainGo := findNodeByPath(root, "main.go")
	require.NotNil(t, mainGo)
	assert.Equal(t, "main.go", mainGo.Name)
	assert.False(t, mainGo.IsFolder)

	// Nested file
	dbGo := findNodeByPath(root, "internal", "pkg", "database", "database.go")
	require.NotNil(t, dbGo)
	assert.Equal(t, "database.go", dbGo.Name)
}

func findNodeByPath(n *Node, path ...string) *Node {
	current := n
	for _, name := range path {
		found := false
		for _, child := range current.Content {
			if child.Name == name {
				current = child
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	return current
}

func TestParseLine_EmptyLine(t *testing.T) {
	t.Parallel()
	root := &Node{}
	root.ParseLine("")
	assert.Nil(t, root.Content)
	assert.Equal(t, "", root.Name)
}

func TestParseLine_RootNode(t *testing.T) {
	t.Parallel()
	root := &Node{}
	root.ParseLine("root")

	assert.Equal(t, "root", root.Name)
	assert.Equal(t, 0, root.Level)
	assert.False(t, root.IsFolder)
	assert.Nil(t, root.ParentNode)
	assert.Nil(t, root.Content)
}

func TestParseLine_FolderNode(t *testing.T) {
	t.Parallel()
	root := &Node{}
	root.ParseLine("  subfolder/")

	assert.Equal(t, "subfolder", root.Content[0].Name)
	assert.Equal(t, 1, root.Content[0].Level)
	assert.True(t, root.Content[0].IsFolder)
}

func TestParseLine_MultipleLevels(t *testing.T) {
	t.Parallel()
	root := &Node{}
	lines := []string{
		"root/",
		"  sub1/",
		"    file1",
		"  sub2/",
		"    file2",
	}
	for _, line := range lines {
		root.ParseLine(line)
	}

	assert.Equal(t, "root", root.Name)
	assert.Equal(t, 2, len(root.Content)) // sub1, sub2

	sub1 := root.Content[0]
	assert.Equal(t, "sub1", sub1.Name)
	assert.True(t, sub1.IsFolder)
	assert.Equal(t, 1, sub1.Level)
	assert.Equal(t, 1, len(sub1.Content))
	assert.Equal(t, "file1", sub1.Content[0].Name)

	sub2 := root.Content[1]
	assert.Equal(t, "sub2", sub2.Name)
	assert.True(t, sub2.IsFolder)
	assert.Equal(t, 1, sub2.Level)
	assert.Equal(t, 1, len(sub2.Content))
	assert.Equal(t, "file2", sub2.Content[0].Name)
}

func TestParseLine_SiblingFolders(t *testing.T) {
	t.Parallel()
	root := &Node{}
	root.ParseLine("root/")
	root.ParseLine("  folderA/")
	root.ParseLine("  folderB/")

	assert.Equal(t, 2, len(root.Content))
	assert.Equal(t, "folderA", root.Content[0].Name)
	assert.Equal(t, "folderB", root.Content[1].Name)
}

func TestParseLine_MixedFilesAndFolders(t *testing.T) {
	t.Parallel()
	root := &Node{}
	root.ParseLine("root/")
	root.ParseLine("  folder/")
	root.ParseLine("    fileA")
	root.ParseLine("  fileB")

	assert.Equal(t, 2, len(root.Content))

	folder := root.Content[0]
	fileB := root.Content[1]

	assert.Equal(t, "folder", folder.Name)
	assert.True(t, folder.IsFolder)
	assert.Equal(t, "fileA", folder.Content[0].Name)

	assert.Equal(t, "fileB", fileB.Name)
	assert.False(t, fileB.IsFolder)
}

func TestInsertNode_AsRoot(t *testing.T) {
	t.Parallel()

	parent := &Node{}
	root := &Node{Name: "root", Level: 0, IsFolder: true}

	parent.insertNode(root)

	assert.Equal(t, "root", parent.Name)
	assert.True(t, parent.IsFolder)
	assert.Equal(t, 0, parent.Level)
	assert.Nil(t, parent.ParentNode)
	assert.Nil(t, parent.Content)
}

func TestInsertNode_AsDirectChild(t *testing.T) {
	t.Parallel()

	parent := &Node{Name: "root", Level: 0, IsFolder: true}
	child := &Node{Name: "child", Level: 1}

	parent.insertNode(child)

	assert.Equal(t, 1, len(parent.Content))
	assert.Equal(t, "child", parent.Content[0].Name)
	assert.Equal(t, parent, parent.Content[0].ParentNode)
}

func TestInsertNode_AsNestedChild(t *testing.T) {
	t.Parallel()

	root := &Node{Name: "root", Level: 0, IsFolder: true}
	child1 := &Node{Name: "folderA", Level: 1, IsFolder: true}
	child2 := &Node{Name: "fileB", Level: 2, IsFolder: false}

	root.insertNode(child1)
	root.insertNode(child2) // should be inserted under child1

	assert.Equal(t, 1, len(root.Content))
	assert.Equal(t, "folderA", root.Content[0].Name)

	folderA := root.Content[0]
	assert.Equal(t, 1, len(folderA.Content))
	assert.Equal(t, "fileB", folderA.Content[0].Name)
	assert.Equal(t, folderA, folderA.Content[0].ParentNode)
}

func TestIsRoot(t *testing.T) {
	t.Parallel()

	root := &Node{Level: 0}
	nonRoot := &Node{Level: 1}

	assert.True(t, root.isRoot())
	assert.False(t, nonRoot.isRoot())
}

func TestIsParentOf(t *testing.T) {
	t.Parallel()

	parent := &Node{Level: 2}
	child := &Node{Level: 3}
	notChild := &Node{Level: 4}

	assert.True(t, parent.isParentOf(child))
	assert.False(t, parent.isParentOf(notChild))
}

func TestIsLevelHigherThan(t *testing.T) {
	t.Parallel()

	high := &Node{Level: 1}
	low := &Node{Level: 3}

	assert.True(t, high.isLevelHigherThan(low))
	assert.False(t, low.isLevelHigherThan(high))
}
