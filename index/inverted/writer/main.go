package main

import (
	"fmt"
	"os"

	"github.com/lamzin/search-engine/index/inverted/writer/lib"
)

func main() {
	if len(os.Args) != 3 {
		panic("usage: cmd path/to/articles path/to/index")
	}
	var articlesPath, indexPath = os.Args[1], os.Args[2]
	fmt.Printf("Articles root: %s\n", articlesPath)
	fmt.Printf("Index root: %s\n", indexPath)

	writer, err := lib.NewInvertedIndexWriter(articlesPath, indexPath)
	if err != nil {
		panic(err.Error())
	}
	defer writer.Close()
	if err := writer.Run(); err != nil {
		panic(err.Error())
	}
	fmt.Println("Index built!")
}
