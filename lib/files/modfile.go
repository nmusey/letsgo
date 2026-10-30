package files

import (
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

// Returns the module name of the nearest go.mod file upwards in the filetree
func ModuleName() (string, error) {
	modfilePath, err := findNearestModfile()
	if err != nil {
		return "", err
	}

	contents, err := os.ReadFile(modfilePath)
	if err != nil {
		return "", err
	}

	moduleName := modfile.ModulePath(contents)
	if moduleName == "" {
		return "", errors.New("no module declaration found in " + modfilePath)
	}

	return moduleName, nil
}

// Traverses up the path to find a go.mod file
func findNearestModfile() (string, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		candidate := filepath.Join(currentDirectory, "go.mod")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}

		parent := filepath.Dir(currentDirectory)
		if parent == currentDirectory {
			return "", errors.New("no go.mod found in the current directory or any parent")
		}
		currentDirectory = parent
	}
}
