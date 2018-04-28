package model

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// NewIndexData create IndexData struct
func NewIndexData(dataRoot string) *IndexData {
	return &IndexData{
		DataRoot: dataRoot,
	}
}

// IndexData struct
type IndexData struct {
	DataRoot string

	// files
	filePaths []string
}

// DumpDocument will dump document
func (data *IndexData) DumpDocument(documentName string, documentBody string) error {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	fileName := reg.ReplaceAllString(strings.ToLower(documentName), "")
	if len(fileName) > 64 {
		fileName = fileName[:64]
	}

	fileFolder := "_short"
	if len(fileName) > 2 {
		fileFolder = strings.ToLower(fileName[:2])
	}
	os.Mkdir(filepath.Join(data.DataRoot, fileFolder), os.ModePerm)

	filePath := filepath.Join(fileFolder, fileName+".txt")
	data.filePaths = append(data.filePaths, filePath)

	filePath = filepath.Join(data.DataRoot, filePath)

	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		return err
	}
	if _, err := file.WriteString(documentBody); err != nil {
		return err
	}
	return nil
}

// Close IndexData
func (data *IndexData) Close() error {
	file, err := os.Create(data.docsFilePath())
	if err != nil {
		return err
	}
	for _, path := range data.filePaths {
		if _, err := file.WriteString(path + "\n"); err != nil {
			return err
		}
	}
	return nil
}

// GetDocsPaths return document's paths
func (data *IndexData) GetDocsPaths() ([]string, error) {
	file, err := os.Open(data.docsFilePath())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (data *IndexData) docsFilePath() string {
	return filepath.Join(data.DataRoot, "docs.txt")
}
