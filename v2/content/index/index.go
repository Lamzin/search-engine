package main

import (
	"fmt"
	"os"

	"github.com/lamzin/search-engine/v2/index/builder"
)

var articlesPath, indexPath string

func main() {
	if len(os.Args) != 3 {
		panic("usage: cmd path/to/articles path/to/index")
	}
	articlesPath, indexPath = os.Args[1], os.Args[2]

	fmt.Printf("Articles root: %s\n", articlesPath)
	fmt.Printf("Index root: %s\n", indexPath)

	m()
	r()
}

func m() {
	mapper := builder.NewIndexBuilderMapper(articlesPath, indexPath)
	if err := mapper.Run(); err != nil {
		panic(err)
	}
}

func r() {
	reducer := builder.NewIndexBuilderReducer(indexPath)
	if err := reducer.Run(); err != nil {
		panic(err)
	}
}
