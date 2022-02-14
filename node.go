package gofind

import (
	"os"
)

type NodePath interface {
	Path() string
	IsDir() bool
	Type() os.FileMode
	Nodes() []NodePath
	Append(node NodePath)
	Name() string
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
	return n.info.IsDir()
}

func (n *node) Name() string {
	return n.info.Name()
}

func (n *node) Nodes() []NodePath {
	return n.nodes
}

func (n *node) Append(node NodePath) {
	if !n.info.IsDir() {
		return
	}

	n.nodes = append(n.nodes, node)
}

func (n *node) Type() os.FileMode {
	return n.info.Mode().Type()
}
