package main

import (
	"fmt"
	"os"

	"github.com/lamzin/search-engine/algos/lexeme"
	"github.com/lamzin/search-engine/doc"
	"github.com/lamzin/search-engine/mysql/index"
)

var articlesPath, indexPath string

func main() {
	if len(os.Args) != 3 {
		panic("usage: cmd path/to/articles path/to/index")
	}
	articlesPath, indexPath = os.Args[1], os.Args[2]

	fmt.Printf("Articles root: %s\n", articlesPath)
	fmt.Printf("Index root: %s\n", indexPath)

	index, err := sqlindex.NewSQLIndex(indexPath)
	if err != nil {
		panic(err)
	}

	iterator := doc.NewDocAllIterator(articlesPath)
	var lexemizer lexeme.Parser
	var lexemes int
	for i := 0; iterator.Scan() && i < 1000; i++ {
		d := iterator.Doc()
		wordFrequencies := lexemizer.Parse(d.String())
		for _, wordFrequency := range wordFrequencies {
			if err = index.Add(wordFrequency.Word, d.ID, wordFrequency.Frequency); err != nil {
				panic(err)
			}
			lexemes++
			fmt.Printf("\rdocs: %d, lexemes: %d", i, lexemes)
		}
	}

	err = index.Close()
	if err != nil {
		panic(err)
	}

}
