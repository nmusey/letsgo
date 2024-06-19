package utils

import (
	"strings"
)

func ReplaceAllInString(str string, replacements map[string]string) string {
	replacementArray := []string{}
	for key, value := range replacements {
		replacementArray = append(replacementArray, key, value)
	}
	replacer := strings.NewReplacer(replacementArray...)
	return replacer.Replace(str)
}
