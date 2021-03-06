package gofind

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Find interface {
	Run()
	Name() string
	Results() []string
	Node() NodePath
}

type find struct {
	expression string
	fileSystem fs.FS
	node       NodePath
	mode       *fs.FileMode
	name       *string
}

func NewFind(expression string, fileSystem fs.FS, mode *fs.FileMode, name *string) *find {
	return &find{
		expression: expression,
		fileSystem: fileSystem,
		mode:       mode,
		name:       name,
	}
}

func (f *find) Run() {
	dirPath := filepath.Dir(f.expression)
	dir, err := f.fileSystem.Open(dirPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dir.Close()

	info, err := dir.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	f.node = NewNode(dirPath, info)

	match, err := filepath.Match(fmt.Sprintf("%s", dirPath), filepath.Clean(f.expression))
	if err != nil {
		fmt.Println(err)
	}

	path := fmt.Sprintf("%s", f.expression)
	if match {
		path = fmt.Sprintf("%s*", f.expression)
	}
	readDir(f.node, path, f.fileSystem)
}

func (f *find) Results() []string {
	var result []string
	if f.Node().Nodes() != nil {
		for _, node := range f.Node().Nodes() {
			readNodes(node, &result, f.mode, f.name)
		}
	}

	return result
}

func (f *find) Node() NodePath {
	return f.node
}

func readNodes(node NodePath, result *[]string, mode *os.FileMode, name *string) {
	if node.Nodes() != nil {
		for _, n := range node.Nodes() {
			readNodes(n, result, mode, name)
		}
	}

	if mode != nil && *mode != node.Type() {
		return
	}

	if name != nil && *name != node.Name() {
		return
	}

	*result = append(*result, node.Path())
}

func readDir(node NodePath, expression string, fileSystem fs.FS) {
	dirEntries, err := fs.ReadDir(fileSystem, node.Path())
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	for _, value := range dirEntries {
		info, err := value.Info()
		if err != nil {
			fmt.Printf("error: %v", err)
		}

		newNode := NewNode(filepath.Join(node.Path(), info.Name()), info)

		match, err := filepath.Match(expression, newNode.Path())
		if err != nil {
			fmt.Printf("error: %v", err)
		}

		if match {
			node.Append(newNode)
			if newNode.IsDir() {
				readDir(newNode, fmt.Sprintf("%s/*", newNode.Path()), fileSystem)
			}
		}
	}
}
