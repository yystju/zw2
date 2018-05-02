package zw2

import (
	"os"
	"strings"
)

type ZipFileNode struct {
	Parent   *ZipFileNode
	Children []*ZipFileNode

	Name    string
	Payload []byte
	Info    os.FileInfo
}

func (n *ZipFileNode) String() string {
	return n.Name
}

type TreeVisiter func(int, *ZipFileNode) bool

func (n *ZipFileNode) Walk(indent int, visitor TreeVisiter) bool {
	r := visitor(indent, n)

	if r {
		for _, child := range n.Children {
			r = r && child.Walk(indent+1, visitor)
		}
	}

	return r
}

func (n *ZipFileNode) GetFullPath() []string {
	parent := n

	r := make([]string, 0)

	for parent != nil {
		r = append(r, parent.Name)
		parent = parent.Parent
	}

	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	return r
}

func (n *ZipFileNode) FindChildByName(name string) *ZipFileNode {
	for _, child := range n.Children {
		if child.Name == name {
			return child
		}
	}

	return nil
}

func (n *ZipFileNode) FindChildByPathName(pathname string) *ZipFileNode {
	if strings.HasPrefix(pathname, "/") {
		pathname = pathname[1:]
	}

	if strings.HasSuffix(pathname, "/") {
		pathname = pathname[:len(pathname)-1]
	}

	paths := strings.Split(pathname, "/")

	parent := n

	for _, p := range paths {
		if parent == nil {
			break
		}

		found := false

		for _, child := range parent.Children {
			if child.Name == p {
				found = true

				parent = child
				continue
			}
		}

		if !found {
			parent = nil
		}
	}

	return parent
}

func (n *ZipFileNode) AddChild(name string, payload []byte, info os.FileInfo) *ZipFileNode {
	child := n.FindChildByName(name)

	if child == nil {
		child = NewZipFileNode(n, name, payload, info)

		n.Children = append(n.Children, child)
	}

	return child
}

func (n *ZipFileNode) AddDescendants(pathname string, payload []byte, info os.FileInfo) {
	if strings.HasSuffix(pathname, "/") {
		pathname = pathname[:len(pathname)-1]
	}

	paths := strings.Split(pathname, "/")

	parent := n

	for i, p := range paths {
		if i == (len(paths) - 1) {
			parent = parent.AddChild(p, payload, info)
		} else {
			parent = parent.AddChild(p, nil, nil)
		}
	}
}

func NewZipFileNode(parent *ZipFileNode, name string, payload []byte, info os.FileInfo) *ZipFileNode {
	r := new(ZipFileNode)

	r.Parent = parent
	r.Name = name
	r.Payload = payload
	r.Info = info

	return r
}
