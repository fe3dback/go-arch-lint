package pathsort

import (
	"sort"
	"strings"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type (
	node struct {
		name   string
		ptr    map[string]*node
		child  []*node
		values []*models.FileDescriptor
	}
)

func newNode(name string) *node {
	return &node{
		name:   name,
		ptr:    make(map[string]*node, 3),
		child:  make([]*node, 0, 3),
		values: make([]*models.FileDescriptor, 0, 2),
	}
}

func (t *node) append(value *models.FileDescriptor) {
	parts := strings.Split(string(value.PathRel), "/")

	parent := t
	length := len(parts)

	for ind, part := range parts {
		isLast := ind == length-1
		var cursor *node

		if _, exist := parent.ptr[part]; !exist {
			node := newNode(part)
			parent.ptr[part] = node
			parent.child = append(parent.child, node)
			cursor = node
		} else {
			cursor = parent.ptr[part]
		}

		if isLast {
			cursor.values = append(cursor.values, value)
		}

		parent = cursor
	}
}

func (t *node) sortLevels() {
	t.ptr = nil

	sort.Slice(t.values, func(i, j int) bool {
		if t.values[i].IsDir != t.values[j].IsDir {
			return t.values[i].IsDir
		}

		return t.values[i].PathRel <= t.values[j].PathRel
	})

	sort.Slice(t.child, func(i, j int) bool {
		isDirI := len(t.child[i].child) > 0
		isDirJ := len(t.child[j].child) > 0

		if isDirI != isDirJ {
			return isDirI
		}

		return t.child[i].name <= t.child[j].name
	})

	for _, child := range t.child {
		child.sortLevels()
	}
}

func (t *node) traversalDepthFirst() []*models.FileDescriptor {
	return recursiveExtractLeafsDepthFirst(t)
}

func recursiveExtractLeafsDepthFirst(node *node) []*models.FileDescriptor {
	list := make([]*models.FileDescriptor, 0, len(node.values)+len(node.ptr))

	list = append(list, node.values...)

	for _, child := range node.child {
		list = append(list, recursiveExtractLeafsDepthFirst(child)...)
	}

	return list
}
