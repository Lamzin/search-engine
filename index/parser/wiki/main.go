package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lamzin/search-engine/index/model/doc"
)

func main() {
	if len(os.Args) != 3 {
		panic("usage: cmd path/to/articles path/to/engine/root")
	}
	var articles, engineRoot = os.Args[1], os.Args[2]
	fmt.Printf("Articles: %s\n", articles)
	fmt.Printf("Engine root: %s\n", engineRoot)

	docManager, err := doc.NewDocManager(engineRoot)
	if err != nil {
		panic(err.Error())
	}
	defer docManager.Close()

	splitter := WikiArticlesSplitter{
		Articles:   articles,
		docManager: docManager,
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

type Splitter interface {
	Split() error
	DocumentsCount() int
	LinesCount() int
}

type WikiArticlesSplitter struct {
	Articles   string
	docManager *doc.DocManager

	// stat
	documents int
	lines     int

	// doc
	doc doc.Doc

	// files
	filePaths []string
}

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
			if err := s.docManager.DumpDocument(s.doc); err != nil {
				fmt.Println(err.Error())
			}
			s.documents++
			s.doc = doc.Doc{
				DocInfo: doc.DocInfo{
					Name: strings.TrimSpace(strings.Replace(line, " =", "", -1)),
				},
				Lines: []string{},
			}
		} else {
			s.doc.AddLine(line)
		}

		if s.lines%10000 == 0 {
			fmt.Printf("\rLines: %dM, docs: %d", s.lines/1000000, s.documents)
		}

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
