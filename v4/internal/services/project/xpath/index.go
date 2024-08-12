package xpath

import (
	"path/filepath"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type index struct {
	full           []models.FileDescriptor
	files          map[models.PathRelative]*models.FileDescriptor
	directories    map[models.PathRelative]*models.FileDescriptor
	directoryFiles map[models.PathRelative][]*models.FileDescriptor
}

func newIndex() *index {
	return &index{
		full:           make([]models.FileDescriptor, 0, 256),
		files:          make(map[models.PathRelative]*models.FileDescriptor, 256),
		directories:    make(map[models.PathRelative]*models.FileDescriptor, 64),
		directoryFiles: make(map[models.PathRelative][]*models.FileDescriptor, 64),
	}
}

func (ind *index) appendToIndex(path models.PathRelative, src models.FileDescriptor) {
	parent := models.PathRelative(filepath.Dir(string(path)))

	ind.full = append(ind.full, src)
	descriptor := &ind.full[len(ind.full)-1]

	// add file to index
	if !descriptor.IsDir {
		ind.files[path] = descriptor
	}

	// create dirs index if not exist
	switch descriptor.IsDir {
	case true:
		ind.directories[path] = descriptor
		if _, exists := ind.directoryFiles[path]; !exists {
			ind.directoryFiles[path] = make([]*models.FileDescriptor, 0, 8)
		}
	case false:
		if _, exists := ind.directoryFiles[parent]; !exists {
			ind.directoryFiles[parent] = make([]*models.FileDescriptor, 0, 8)
		}
	}

	// add file to dir index
	if !descriptor.IsDir {
		ind.directoryFiles[parent] = append(ind.directoryFiles[parent], descriptor)
	}
}

func (ind *index) fileAt(path models.PathRelative) (models.FileDescriptor, bool) {
	dst, ok := ind.files[path]
	if !ok {
		return models.FileDescriptor{}, false
	}

	if dst.IsDir {
		return models.FileDescriptor{}, false
	}

	return *dst, true
}

func (ind *index) directoryAt(path models.PathRelative) (models.FileDescriptor, bool) {
	dst, ok := ind.directories[path]
	if !ok {
		return models.FileDescriptor{}, false
	}

	if !dst.IsDir {
		return models.FileDescriptor{}, false
	}

	return *dst, true
}

func (ind *index) each(fn func(models.FileDescriptor)) {
	for _, descriptor := range ind.full {
		fn(descriptor)
	}
}
