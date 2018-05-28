package main

import (
	"fmt"

	"github.com/lamzin/search-engine/algos/compressor/numbers"
)

func main() {
	test([]int{1})
	test([]int{1, 1, 1, 1, 1})
	test([]int{1, 1, 1, 1, 1, 2, 3, 4})
	test([]int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024})
	test([]int{})
	test([]int{4294967295})
	test([]int{2, 2, 2, 2, 2, 2, 2, 2})
	test([]int{2, 2, 2, 2, 2, 2, 2, 2, 2})
	test([]int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2})
	test([]int{8, 8, 8, 8, 8, 8, 8, 8})

	test2([]int{1}, []int{1})
	test2([]int{1, 2, 3, 4, 5, 6}, []int{7, 8, 9, 10, 11, 12})
	test2([]int{1, 1, 1, 1, 1, 2, 3, 4}, []int{1, 1, 1, 1, 1, 2, 3, 4})
}

func test(numbers []int) {
	elias := numberscompressor.EliasGammaCodes{}
	bytes, payload, _ := elias.Compress(numbers, 0, 0)

	numbersEncoded, _ := elias.Decompress(bytes)
	fmt.Println(numbers, "-->", bytes, payload, len(bytes)*8, "-->", numbersEncoded)
	if len(numbers) != len(numbersEncoded) {
		panic("diff len")
	}
	for i := 0; i < len(numbers); i++ {
		if numbers[i] != numbersEncoded[i] {
			panic(fmt.Errorf("diff numbers on position %d", i))
		}
	}
}

func test2(numbers1 []int, numbers2 []int) {
	elias := numberscompressor.EliasGammaCodes{}
	bytes1, payload1, _ := elias.Compress(numbers1, 0, 0)
	bytes2, payload2, _ := elias.Compress(numbers2, bytes1[len(bytes1)-1], payload1)
	bytes := append(bytes1[:len(bytes1)-1], bytes2...)

	numbers := append(numbers1, numbers2...)

	numbersEncoded, _ := elias.Decompress(bytes)
	fmt.Println(numbers1, numbers2, "-->", bytes, len(bytes1)*8-8+payload2, len(bytes)*8, "-->", numbersEncoded)
	if len(numbers) != len(numbersEncoded) {
		panic("diff len")
	}
	for i := 0; i < len(numbers); i++ {
		if numbers[i] != numbersEncoded[i] {
			panic(fmt.Errorf("diff numbers on position %d", i))
		}
	}
}
