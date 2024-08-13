package xpath_test

import (
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type filePathStub struct {
	path     models.PathAbsolute
	isDir    bool
	expected bool
}

type fileScannerTestStub struct {
	content []filePathStub
}

func newFileScannerTestStub(content []filePathStub) *fileScannerTestStub {
	return &fileScannerTestStub{
		content: content,
	}
}

func (s *fileScannerTestStub) Scan(_ string, fn func(path string, isDir bool) error) error {
	for _, stub := range s.content {
		err := fn(string(stub.path), stub.isDir)
		if err != nil {
			return err
		}
	}

	return nil
}
