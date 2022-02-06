package gofind

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestNewNode(t *testing.T) {

	var testData = []struct {
		fs    fstest.MapFS
		name  string
		isDir bool
		path  string
	}{
		{
			fs: fstest.MapFS{
				"test-dir": &fstest.MapFile{
					Mode: os.ModeDir,
				},
			},
			name:  "test-dir",
			isDir: true,
			path:  "/some-path/testing/test-dir",
		},
		{
			fs: fstest.MapFS{
				"test-file": &fstest.MapFile{},
			},
			name:  "test-file",
			isDir: false,
			path:  "/some-path/testing/test-file",
		},
	}

	for _, value := range testData {
		info, _ := value.fs.Stat(value.name)
		node := NewNode(value.path, info)

		if node.IsDir() != value.isDir {
			t.Error("file mode is not correct")
		}

		if node.Path() != value.path {
			t.Errorf("path %s != %s ", value.path, node.Path())
		}
	}
}

func TestAppendChildNodesToDirNode(t *testing.T) {

	fs := fstest.MapFS{
		"test-dir": &fstest.MapFile{
			Mode: os.ModeDir,
		},
		"test-file": &fstest.MapFile{},
	}

	dir, _ := fs.Stat("test-dir")
	node := NewNode(filepath.Join("/some-path", dir.Name()), dir)

	fileInfo, _ := fs.Stat("test-file")

	node.Append(NewNode(filepath.Join(node.Path(), fileInfo.Name()), fileInfo))

	if node.Nodes() == nil || len(node.Nodes()) != 1 {
		t.Error("node should be appended to the nodes slice")
	}

	if node.Nodes()[0].Path() != "/some-path/test-dir/test-file" {
		t.Errorf("child node does not have proper Path: %s", node.Nodes()[0].Path())
	}
}

func TestNoAppendChildNodesToRegulerFileNode(t *testing.T) {
	fs := fstest.MapFS{
		"test-file":  &fstest.MapFile{},
		"other-file": &fstest.MapFile{},
	}

	file, _ := fs.Stat("test-file")
	node := NewNode("/some-path", file)

	otherFile, _ := fs.Stat("other-file")

	node.Append(NewNode(node.Path(), otherFile))

	if node.Nodes() != nil {
		t.Error("nodes should not be appended to the file node")
	}
}
