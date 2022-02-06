package gofind

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type Find interface {
	Run()
	Results() []string
	Node() NodePath
}

type find struct {
	expression string
	fileSystem fs.FS
	node       NodePath
}

func NewFind(expression string, fileSystem fs.FS) *find {
	return &find{
		expression: expression,
		fileSystem: fileSystem,
	}
}

func (f *find) Run() {
	dirPath := filepath.Dir(f.expression)

	match, err := filepath.Match(fmt.Sprintf("%s", dirPath), filepath.Clean(f.expression))
	if err != nil {
		fmt.Println(err)
	}

	f.node = NewNode(dirPath, nil)
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
			readNodes(node, &result)
		}
	}

	return result
}

func (f *find) Node() NodePath {
	return f.node
}

func readNodes(node NodePath, result *[]string) {
	*result = append(*result, node.Path())
	if node.Nodes() != nil {
		for _, n := range node.Nodes() {
			readNodes(n, result)
		}
	}
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
