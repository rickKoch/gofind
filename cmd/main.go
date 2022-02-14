package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/rickKoch/gofind"
)

type file struct{}

func (f file) Open(name string) (fs.File, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return file, nil
}

const usage = `Usage of gofind:
  -t, --type file type
  -n, --name file name
`

func fileType(t string) *fs.FileMode {
	var f fs.FileMode

	switch t {
	case "f":
		return &f
	case "d":
		f = fs.ModeDir
		return &f
	default:
		return nil
	}
}

func main() {

	var t string
	var name string

	flag.StringVar(&t, "type", "", "file type")
	flag.StringVar(&name, "name", "", "file name")

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	expression := flag.Arg(0)

	var n *string
	if name != "" {
		n = &name
	}

	find := gofind.NewFind(expression, file{}, fileType(t), n)
	find.Run()

	if find.Node() == nil {
		fmt.Printf("find: '%s': No such file or direcotry\n", expression)
		return
	}

	for _, v := range find.Results() {
		fmt.Println(v)
	}

}
