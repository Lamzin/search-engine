package utils

import (
	"bufio"
	"io/ioutil"
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

func ReadFileBytes(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.OpenFile(path, os.O_CREATE, 0666)
		defer file.Close()
		return []byte{}, err
	}

	bytes, err := ioutil.ReadFile(path)
	return bytes, err
}

func AppendFile(path string, data []byte) error {
	fileFolder, _ := filepath.Split(path)

	errFolder := os.Mkdir(fileFolder, os.ModePerm)
	if errFolder != nil && !strings.Contains(errFolder.Error(), "file exists") {
		return errFolder
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func WriteFile(path string, data []byte) error {
	fileFolder, _ := filepath.Split(path)

	errFolder := os.Mkdir(fileFolder, os.ModePerm)
	if errFolder != nil && !strings.Contains(errFolder.Error(), "file exists") {
		return errFolder
	}

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}
