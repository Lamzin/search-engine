package main

import (
	"fmt"
	"os"

	"github.com/lamzin/search-engine/algos/lexeme"
	"github.com/lamzin/search-engine/doc"
	"github.com/lamzin/search-engine/v2/index/builder"
)

var (
	lexemizer lexeme.Parser
)

func main() {
	if len(os.Args) != 3 {
		panic("usage: cmd path/to/articles path/to/index")
	}
	var articlesPath, indexPath = os.Args[1], os.Args[2]
	fmt.Printf("Articles root: %s\n", articlesPath)
	fmt.Printf("Index root: %s\n", indexPath)

	iterator := doc.NewDocAllIterator(articlesPath)
	index := builder.NewIndexRAM(indexPath)

	lexemes := 0
	for i := 0; iterator.Scan(); i++ {
		d := iterator.Doc()
		wordFrequencies := lexemizer.Parse(d.String())
		for _, wordFrequency := range wordFrequencies {
			if index.CanAddLexeme() {
				if err := index.AddLexeme(wordFrequency.Word, d.ID, wordFrequency.Frequency); err != nil {
					panic(err)
				}
				lexemes++
				fmt.Printf("\rdocs: %d | lexemes: %d", i, lexemes)
			} else {
				fmt.Println("start index dump...")
				if err := index.Dump(); err != nil {
					panic(err)
				}
				fmt.Println("finish index dump")
				index = builder.NewIndexRAM(indexPath)
			}
		}
	}
}
