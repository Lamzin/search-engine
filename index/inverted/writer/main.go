package main

import (
	"fmt"
	"os"

	"github.com/lamzin/search-engine/index/inverted"
	"github.com/lamzin/search-engine/index/inverted/writer/lib"
	"github.com/lamzin/search-engine/index/model/doc"
)

func main() {
	if len(os.Args) != 3 {
		panic("usage: cmd path/to/articles path/to/index")
	}
	var articlesPath, indexPath = os.Args[1], os.Args[2]
	fmt.Printf("Articles root: %s\n", articlesPath)
	fmt.Printf("Index root: %s\n", indexPath)

	iterator := doc.NewDocAllIterator(articlesPath)

	lexemizer := lib.NewLexemizer()

	index := inverted.NewInvertedIndex(indexPath)

	tokensCount := 0
	for i := 0; iterator.Scan(); i++ {
		d := iterator.Doc()

		tokens := lexemizer.Parse(d.String())
		// fmt.Println(d.String())
		// fmt.Println(tokens)
		for _, token := range tokens {
			if err := index.AddToken(&d.DocInfo, token); err != nil {
				fmt.Println("error adding token:", err)
				return
			}
			tokensCount++
		}
		fmt.Printf("\rdocs: %d, tokens: %d", i, tokensCount)
		if i == 20000 {
			break
		}
	}

	fmt.Printf("lexemizer stat: all - %.3f -> unique %.3f -> words only %.3f -> long words %.3f -> stem %.3f\n",
		float64(lexemizer.StatAll)/float64(lexemizer.StatAll),
		float64(lexemizer.StatUnique)/float64(lexemizer.StatAll),
		float64(lexemizer.StatWordsOnly)/float64(lexemizer.StatAll),
		float64(lexemizer.StatLongWords)/float64(lexemizer.StatAll),
		float64(lexemizer.StatStem)/float64(lexemizer.StatAll))

	fmt.Println("iteration error:", iterator.Err())
}
