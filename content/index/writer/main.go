package main

import (
	"fmt"
	"os"

	"github.com/lamzin/search-engine/algos/lexeme"
	"github.com/lamzin/search-engine/doc"
	"github.com/lamzin/search-engine/index/writer"
)

func main() {
	if len(os.Args) != 3 {
		panic("usage: cmd path/to/articles path/to/index")
	}
	var articlesPath, indexPath = os.Args[1], os.Args[2]
	fmt.Printf("Articles root: %s\n", articlesPath)
	fmt.Printf("Index root: %s\n", indexPath)

	iterator := doc.NewDocAllIterator(articlesPath)

	lexemizer := lexeme.NewParser()

	// var index indexwriter.Writer = indexwriter.NewMultiLexemePerFile(indexPath)
	var index indexwriter.Writer = indexwriter.NewIndexDBWriter(indexPath)

	tokenStat := make(map[string]int, 0)

	tokensCount := 0
	for i := 0; iterator.Scan(); i++ {
		d := iterator.Doc()

		tokens := lexemizer.Parse(d.String())
		// fmt.Println(d.String())
		// fmt.Println(tokens)
		for _, token := range tokens {
			tokenStat[token]++
			if err := index.AddLexeme(d.DocInfo.ID, token); err != nil {
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

	leng := 0
	for k, v := range tokenStat {
		fmt.Println(k, v)
		leng += len(k)
	}

	fmt.Println("unique tokens:", len(tokenStat))
	fmt.Println("tokens:", tokensCount)
	fmt.Println("average tokens:", tokensCount/len(tokenStat))
	fmt.Println("length tokens:", leng)

	fmt.Printf("lexemizer stat: all - %.3f -> unique %.3f -> words only %.3f -> long words %.3f -> stem %.3f\n",
		float64(lexemizer.StatAll)/float64(lexemizer.StatAll),
		float64(lexemizer.StatUnique)/float64(lexemizer.StatAll),
		float64(lexemizer.StatWordsOnly)/float64(lexemizer.StatAll),
		float64(lexemizer.StatLongWords)/float64(lexemizer.StatAll),
		float64(lexemizer.StatStem)/float64(lexemizer.StatAll))

	fmt.Println("iteration error:", iterator.Err())
}
