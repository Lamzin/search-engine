package main

import (
	"fmt"
	"math/rand"

	"github.com/lamzin/search-engine/algos/compressor/numbers"
)

func main() {
	comp := numberscompressor.VariableByteCodes{}

	for i := 0; i < 1000000; i++ {
		fmt.Printf("\r%d", i)
		var arr []int
		for j := 0; j < (int)(rand.Int31n(100000)); j++ {
			arr = append(arr, (int)(rand.Int31n(1<<30)))
		}
		bytes, err := comp.Compress(arr)
		if err != nil {
			panic(err)
		}
		b, err := comp.Decompress(bytes)
		if err != nil {
			panic(err)
		}
		if len(arr) != len(b) {
			panic(fmt.Sprintf("%s %s", arr, b))
		}
		for j := 0; j < len(arr); j++ {
			if arr[j] != b[j] {
				panic(fmt.Sprintf("%s %s", arr, b))
			}
		}
	}
}
