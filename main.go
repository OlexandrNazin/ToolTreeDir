package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"toolTreeDir/sizeByteConvert"
)

type Node interface {
	fmt.Stringer
}

type Directory struct {
	name    string
	folders []Node
}

type File struct {
	name    string
	size    int64
	modTime string
}

func main() {
	stdOut := os.Stdout
	fmt.Fprintf(stdOut, "%s\n", os.Args[0])
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		fmt.Println(`	
	Enter the relative directory path.
		
		--help, more info.
`)
		return
	} else if os.Args[1] == `--help` {

		fmt.Println(`	Displays the folder tree of the specified directory. 

	Specify the full path, or use "." for current.

 	-f flag will display information about the nested files (size in bytes, last modified date).
`)
		return
	} else {
		path := os.Args[1]

		printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
		err := dirTree(stdOut, path, printFiles)
		if err != nil {
			fmt.Printf("Open %v: no such file or directory\n", os.Args[1])
			return
		}
	}
}

func (file *File) String() string {
	if file.size == 0 {
		return file.name + " (empty) " + file.modTime[:19]
	}
	return fmt.Sprintf("%v (%+v), %+v", file.name, sizeByteConvert.Convert(file.size), file.modTime[:19])
}

func (directory *Directory) String() string {
	return directory.name
}

func readDir(path string, nodes []Node, withFiles bool) (error, []Node) {
	files, err := ioutil.ReadDir(path)
	for _, info := range files {
		if !(info.IsDir() || withFiles) {
			continue
		}

		var newNode Node

		if info.IsDir() {
			_, folders := readDir(filepath.Join(path, info.Name()), []Node{}, withFiles)
			newNode = &Directory{info.Name(), folders}
		} else {
			newNode = &File{info.Name(), info.Size(), info.ModTime().String()}
		}

		nodes = append(nodes, newNode)
	}
	return err, nodes
}

func printDir(out io.Writer, nodes []Node, prefixes []string) {
	if len(nodes) == 0 {
		return
	}

	fmt.Fprintf(out, "%s", strings.Join(prefixes, ""))

	node := nodes[0]

	if len(nodes) == 1 {
		fmt.Fprintf(out, "%s%s\n", "└───", node)
		if directory, ok := node.(*Directory); ok {
			printDir(out, directory.folders, append(prefixes, "    "))
		}
		return
	}

	fmt.Fprintf(out, "%s%s\n", "├───", node)
	if directory, ok := node.(*Directory); ok {
		printDir(out, directory.folders, append(prefixes, "│    "))
	}

	printDir(out, nodes[1:], prefixes)
}

func dirTree(out io.Writer, path string, pritnFiles bool) error {
	err, nodes := readDir(path, []Node{}, pritnFiles)
	printDir(out, nodes, []string{})
	return err
}
