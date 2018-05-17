package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	n, err := file.WriteAt([]byte("Hello!"), 100)
	fmt.Println(n, err)

	err = os.Truncate("test.txt", 1000)
	fmt.Println(err)
}
