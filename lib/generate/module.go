package generate

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ModuleName(dir string) (string, error) {
	path := filepath.Join(dir, "go.mod")

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("generate: no go.mod found in %s — run this command from your project root", dir)
		}
		return "", fmt.Errorf("generate: reading go.mod: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if after, ok := strings.CutPrefix(line, "module "); ok {
			return strings.TrimSpace(after), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("generate: reading go.mod: %w", err)
	}

	return "", fmt.Errorf("generate: go.mod in %s has no module directive", dir)
}
