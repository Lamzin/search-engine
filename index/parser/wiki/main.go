package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lamzin/search-engine/index/model"
)

func main() {
	if len(os.Args) != 3 {
		panic("usage: cmd path/to/articles path/to/engine/root")
	}
	var articles, engineRoot = os.Args[1], os.Args[2]
	fmt.Printf("Articles: %s\n", articles)
	fmt.Printf("Engine root: %s\n", engineRoot)

	indexData := model.NewIndexData(engineRoot)
	defer indexData.Close()

	splitter := WikiArticlesSplitter{
		Articles:  articles,
		IndexData: indexData,
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
	Articles  string
	IndexData *model.IndexData

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
			if err := s.IndexData.DumpDocument(s.curDocumentName, s.curDocument); err != nil {
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
