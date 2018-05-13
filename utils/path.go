package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

var fileNameReg *regexp.Regexp

const (
	shortFolderName = "_short"
	fileExtention   = ".txt"

	maxFileNameLen = 64
	fileFolderLen  = 2
)

func init() {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err.Error())
	}
	fileNameReg = reg
}

func FilePath(s string) (string, string) {
	fileName := fileNameReg.ReplaceAllString(strings.ToLower(s), "")
	if len(fileName) > maxFileNameLen {
		fileName = fileName[:maxFileNameLen]
	}

	fileFolder := shortFolderName
	if len(fileName) > fileFolderLen {
		fileFolder = strings.ToLower(fileName[:fileFolderLen])
	}
	return fileFolder, filepath.Join(fileFolder, fileName+fileExtention)
}
