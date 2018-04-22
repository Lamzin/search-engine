package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		panic("usage: cmd path/to/articles path/to/engine/root")
	}
	var articles, engineRoot = os.Args[1], os.Args[2]
	fmt.Printf("Articles: %s\n", articles)
	fmt.Printf("Engine root: %s\n", engineRoot)

	splitter := WikiArticlesSplitter{
		Articles:   articles,
		EngineRoot: engineRoot,
	}

	if err := split(&splitter); err != nil {
		log.Fatal(err.Error())
	}
}

func split(splitter Splitter) error {
	if err := splitter.Split(); err != nil {
		return err
	}
	return nil
}

// Splitter interface
type Splitter interface {
	Split() error
	DocumentsCount() int
	LinesCount() int
}

// WikiArticlesSplitter interface
type WikiArticlesSplitter struct {
	Articles   string
	EngineRoot string

	// stat
	documents int
	lines     int

	// docs
	curDocumentName string
	curDocument     string

	// files
	filePaths []string
}

// Split huge documents into many small
func (s *WikiArticlesSplitter) Split() error {
	file, err := os.Open(s.Articles)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for ; scanner.Scan(); s.lines++ {
		line := scanner.Text()
		if strings.HasPrefix(line, "= ") {
			if err := s.dumpDocument(); err != nil {
				return err
			}
			s.documents++
			s.curDocumentName = strings.TrimSpace(strings.Replace(line, "=", "", -1))
			s.curDocument = s.curDocumentName
		} else {
			s.curDocument += "\n" + line
		}

		if s.lines%10000 == 0 {
			fmt.Printf("\rLines: %dM, docs: %d", s.lines/1000000, s.documents)
		}

		// if s.documents == 100 {
		// 	break
		// }
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := s.dumpFilePaths(); err != nil {
		return err
	}

	return nil
}

// DocumentsCount return documents count
func (s *WikiArticlesSplitter) DocumentsCount() int {
	return s.documents
}

// LinesCount return total lines amount in articles
func (s *WikiArticlesSplitter) LinesCount() int {
	return s.lines
}

func (s *WikiArticlesSplitter) dumpDocument() error {
	fileName := strings.Replace(s.curDocumentName, " ", "_", -1)
	fileName = strings.Replace(fileName, "/", "_", -1)
	if len(fileName) > 64 {
		fileName = fileName[:64]
	}

	fileFolder := "_short"
	if len(fileName) > 2 {
		fileFolder = strings.ToLower(fileName[:2])
	}
	os.Mkdir(filepath.Join(s.EngineRoot, fileFolder), os.ModePerm)

	filePath := filepath.Join(fileFolder, fileName+".txt")
	s.filePaths = append(s.filePaths, filePath)

	filePath = filepath.Join(s.EngineRoot, filePath)

	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		return err
	}
	if _, err := file.WriteString(s.curDocument); err != nil {
		return err
	}
	return nil
}

func (s *WikiArticlesSplitter) dumpFilePaths() error {
	file, err := os.Create(filepath.Join(s.EngineRoot, "docs.txt"))
	if err != nil {
		return err
	}
	for _, path := range s.filePaths {
		if _, err := file.WriteString(path + "\n"); err != nil {
			return err
		}
	}
	return nil
}
