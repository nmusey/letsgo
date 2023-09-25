package utils

import (
	"os"
	"strings"
)

func UpsertFolder(rootpath string) error {
	const mode = 0755
	if _, err := os.Stat(rootpath); os.IsNotExist(err) {
		if err := os.Mkdir(rootpath, mode); err != nil {
			return err
		}
	}

	return nil
}

func UpsertFile(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		file, err := os.Create(filepath)
		if err != nil {
			return err
		}

		defer file.Close()
	}

	return nil
}

func CopyTemplateFile(templatePath string, outpath string, replacements map[string]string) error {
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	filledTemplate := ReplaceAllInString(string(templateBytes), replacements)
	return os.WriteFile(outpath, []byte(filledTemplate), 0644)
}

func ReplaceAllInString(str string, replacements map[string]string) string {
	replacementArray := []string{}
	for key, value := range replacements {
		replacementArray = append(replacementArray, key, value)
	}
	replacer := strings.NewReplacer(replacementArray...)
	return replacer.Replace(str)
}

func CopyFile(filepath string, outpath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = os.WriteFile(outpath, data, 0644)
	return err
}
