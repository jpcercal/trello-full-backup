package internal

import (
	"io/ioutil"
	"os"
	"regexp"
)

var mkdir = os.MkdirAll

// CreateDirectoryRecursively Make a folder if it does not exist already
// (applies it recursively in case the parent folder do not exist)
func CreateDirectoryRecursively(name string) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		if err2 := mkdir(name, os.ModePerm); err2 != nil {
			logger.WithError(err2).WithField("name", name).Fatalf("failed to create directory")
		}
	}
}

// SaveFile Write <content> to the file <filename>
func SaveFile(filename string, content []byte) {
	err := ioutil.WriteFile(filename, content, os.ModePerm)
	if err != nil {
		logger.WithError(err).WithField("filename", filename).Fatalf("failed to save file content")
	}
}

// Sanitize Strip problematic characters on a file name
func Sanitize(name string) string {
	return regexp.MustCompile(`[<>:\/\|\?\*\']`).ReplaceAllString(name, "_")
}
