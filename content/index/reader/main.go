package main

import (
	"fmt"
	"os"

	"github.com/lamzin/search-engine/algos/lexeme"
	"github.com/lamzin/search-engine/doc"
	"github.com/lamzin/search-engine/index/reader"
	"github.com/lamzin/search-engine/index/writer"
)

func main() {
	if len(os.Args) != 3 {
		panic("usage: cmd path/to/articles path/to/index")
	}
	var articlesPath, indexPath = os.Args[1], os.Args[2]
	fmt.Printf("Articles root: %s\n", articlesPath)
	fmt.Printf("Index root: %s\n", indexPath)

	pairs := getPairs(articlesPath)
	write(indexPath, pairs)
	read(indexPath, pairs)
}

type pair struct {
	Doc    int
	Lexeme string
}

func getPairs(docPath string) []pair {
	var pairs []pair
	iterator := doc.NewDocAllIterator(docPath)
	parser := lexeme.NewParser()
	for i := 0; iterator.Scan(); i++ {
		d := iterator.Doc()

		lexemes := parser.Parse(d.String())
		for _, lexeme := range lexemes {
			pairs = append(pairs, pair{Doc: d.ID, Lexeme: lexeme})
		}

		if i == 1000 {
			break
		}
	}
	if err := iterator.Err(); err != nil {
		panic(err)
	}
	return pairs
}

func write(indexPath string, pairs []pair) {
	index := indexwriter.NewIndexDBWriter(indexPath)
	defer index.Close()

	// pairs := []pair{
	// 	{0, "the"},
	// 	{2, "the"},
	// 	{3, "the"},
	// 	{100, "oleh"},
	// 	{200, "iryna"},
	// }

	for _, p := range pairs {
		if err := index.AddLexeme(p.Doc, p.Lexeme); err != nil {
			panic(err)
		}
	}
}

func read(indexPath string, pairs []pair) {
	index, err := indexreader.NewDBReader(indexPath)
	if err != nil {
		panic(err)
	}

	// parser := lexeme.NewParser()
	// words := []string{"the", "oleh", "iryna", "vasya"}

	words := make(map[string][]int)
	for _, p := range pairs {
		words[p.Lexeme] = append(words[p.Lexeme], p.Doc)
	}

	for word, expectedDocIDs := range words {
		// lexemes := parser.Parse(word)
		// fmt.Println(word, "-->", lexemes)
		numbers, err := index.GetDocIDs(word)
		if err != nil {
			panic(err)
		}

		expected := make(map[int]bool)
		actual := make(map[int]bool)
		for _, number := range expectedDocIDs {
			expected[number] = true
		}
		for _, number := range numbers {
			actual[number] = true
		}

		if len(expected) != len(actual) {
			fmt.Println(word, "error: differnt doc arrays", expected, actual)
			panic("diff docs")
		}

		for k, _ := range actual {
			if _, ok := expected[k]; !ok {
				fmt.Println(word, "error: differnt doc arrays", expected, actual)
				panic("diff docs")
			}
		}
		fmt.Println(word, "ok")
	}
}
