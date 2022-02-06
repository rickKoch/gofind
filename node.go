package gofind

import (
	"os"
)

type NodePath interface {
	Path() string
	IsDir() bool
	Nodes() []NodePath
	Append(node NodePath)
}

type node struct {
	path  string
	info  os.FileInfo
	nodes []NodePath
}

func NewNode(path string, info os.FileInfo) NodePath {
	return &node{path: path, info: info}
}

func (n *node) Path() string {
	return n.path
}

func (n *node) IsDir() bool {
	if n.info == nil {
		return true
	}
	return n.info.IsDir()
}

func (n *node) Nodes() []NodePath {
	return n.nodes
}

func (n *node) Append(node NodePath) {
	if !n.IsDir() {
		return
	}

	n.nodes = append(n.nodes, node)
}
