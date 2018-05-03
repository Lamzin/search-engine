package common

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func ReadFile(path string) ([]string, error) {
	inFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func AppendFile(path string, data string) error {
	fileFolder, _ := filepath.Split(path)

	errFolder := os.Mkdir(fileFolder, os.ModePerm)
	if errFolder != nil && !strings.Contains(errFolder.Error(), "file exists") {
		return errFolder
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	return err
}
