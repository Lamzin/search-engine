package main

import (
	"fmt"
	"os"

	"github.com/lamzin/search-engine/index/model"
)

func main() {
	if len(os.Args) != 2 {
		panic("usage: cmd path/to/engine/root")
	}
	var engineRoot = os.Args[1]
	fmt.Printf("Engine root: %s\n", engineRoot)

	indexData := model.NewIndexData(engineRoot)
	fmt.Println(indexData.DataRoot)
}
