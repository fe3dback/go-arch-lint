package path

import (
	"fmt"
	"os"
)

type (
	Type         uint8
	ResolvedPath struct {
		PathType Type
		Paths    []string
	}
)

func ResolvePath(path string) ([]string, error) {
	dirs := make([]string, 0)

	matches, err := glob(path)
	if err != nil {
		return nil, fmt.Errorf("can`t match path mask '%s': %v",
			path,
			err,
		)
	}

	for _, match := range matches {
		fileInfo, err := os.Stat(match)
		if err != nil {
			return nil, fmt.Errorf("nostat '%s': %v", match, err)
		}

		switch mode := fileInfo.Mode(); {
		case mode.IsDir():
			dirs = append(dirs, match)
		default:
			continue
		}
	}

	return dirs, nil
}
