package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// new-book is a soft link to a target. That dir has been excluded through .gitignore so all the content shall be safe.
var srcPath = flag.String("srcPath", "new-book/src", "path to content root where SUMMARY.md shall located")

func main() {
	flag.Parse()

	tree, err := NewNode(filepath.Dir(*srcPath), filepath.Base(*srcPath))
	if err != nil {
		log.Fatal(err)
	}
	tree.Print(0)
}

type Node struct {
	Name     string
	Path     string
	Children []*Node
}

func (n *Node) Print(depth int) {
	prefix := strings.Repeat("  ", depth)
	fmt.Printf("%s[%s](%s)\n", prefix, n.Name, n.Path)
	for _, child := range n.Children {
		child.Print(depth + 1)
	}
}

func NewNode(path, base string) (*Node, error) {
	stat, err := os.Stat(filepath.Join(path, base))
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}
	if stat.IsDir() {
		node, err := NewNodeFromDir(path, base)
		if err != nil {
			return nil, err
		}
		if len(node.Children) == 0 {
			return nil, nil
		}
		return node, err
	}
	return NewNodeFromFile(path, base)
}

func NewNodeFromDir(path, base string) (*Node, error) {
	current := filepath.Join(path, base)

	entries, err := os.ReadDir(current)
	if err != nil {
		return nil, fmt.Errorf("ReadDir: %w", err)
	}

	var children []*Node
	for _, entry := range entries {
		node, err := NewNode(current, entry.Name())
		if err != nil {
			return nil, fmt.Errorf("NewNode from %s: %w", filepath.Join(current, entry.Name()), err)
		}
		// Skip invalid items that got nil,nil from NewNode.
		if node != nil {
			children = append(children, node)
		}
	}

	return &Node{
		Name:     base,
		Path:     filepath.Join(path, base, "README.md"),
		Children: children,
	}, nil
}

func NewNodeFromFile(path, base string) (*Node, error) {
	if filepath.Ext(base) != ".md" || base == "README.md" {
		return nil, nil
	}
	current := filepath.Join(path, base)

	data, err := os.ReadFile(current)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: %w", err)
	}

	return &Node{
		Name:     Digester(data),
		Path:     current,
		Children: nil,
	}, nil
}
