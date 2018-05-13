package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lamzin/search-engine/doc"
)

func main() {
	if len(os.Args) != 3 {
		panic("usage: cmd path/to/articles path/to/engine/root")
	}
	var articles, engineRoot = os.Args[1], os.Args[2]
	fmt.Printf("Articles: %s\n", articles)
	fmt.Printf("Engine root: %s\n", engineRoot)

	docManager := doc.NewDocFileManager(engineRoot)
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
}

type WikiArticlesSplitter struct {
	Articles   string
	docManager doc.DocManager

	// doc
	doc doc.Doc

	// files
	filePaths []string
}

func (s *WikiArticlesSplitter) Split() error {
	var docReader doc.DocReader = doc.NewDocTxtReader(s.Articles)
	defer docReader.Close()

	for count := 0; docReader.Scan(); count++ {
		d := docReader.Doc()
		if err := s.docManager.DumpDocument(d); err != nil {
			return err
		}
		if count%1000 == 0 {
			fmt.Printf("Docs: %d\n", count)
		}
	}

	return docReader.Err()
}
