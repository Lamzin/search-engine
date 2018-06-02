package main

import (
	"fmt"
	"os"

	"github.com/lamzin/search-engine/v2/index/builder"
)

func main() {
	if len(os.Args) != 2 {
		panic("usage: cmd path/to/index")
	}
	var indexPath = os.Args[1]
	fmt.Printf("Index root: %s\n", indexPath)

	index, err := builder.NewIndexRAMStorage(indexPath)
	if err != nil {
		panic(err)
	}

	for _, info := range index.Infos {
		postings, err := index.GetPostings(info.Lexeme)
		if err != nil {
			panic(err)
		}
		frequencies, err := index.GetFrequencies(info.Lexeme)
		if err != nil {
			panic(err)
		}
		if len(postings) != len(frequencies) {
			panic(fmt.Errorf("%s %x %x", info.Lexeme, postings, frequencies))
		}
		fmt.Println(info.Lexeme, postings, frequencies)
	}

}
