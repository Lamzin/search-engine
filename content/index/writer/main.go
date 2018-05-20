package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lamzin/search-engine/algos/lexeme"
	"github.com/lamzin/search-engine/doc"
	"github.com/lamzin/search-engine/index/writer"
)

const (
	queueSize = 1000

	workers = 4

	indexWriterCount = 1
)

var (
	index indexwriter.Writer

	queue chan *doc.Doc = make(chan *doc.Doc, queueSize)
	ready chan bool     = make(chan bool, workers)

	indexes []indexwriter.Writer
)

func main() {
	if len(os.Args) != 3 {
		panic("usage: cmd path/to/articles path/to/index")
	}
	var articlesPath, indexPath = os.Args[1], os.Args[2]
	fmt.Printf("Articles root: %s\n", articlesPath)
	fmt.Printf("Index root: %s\n", indexPath)

	iterator := doc.NewDocAllIterator(articlesPath)

	for i := 0; i < indexWriterCount; i++ {
		path := filepath.Join(indexPath, strconv.Itoa(i))
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			panic(err)
		}
		indexes = append(indexes, indexwriter.NewIndexDBWriter(path))
	}

	for i := 0; i < workers; i++ {
		go worker()
	}

	// tokenStat := make(map[string]int, 0)
	// tokensCount := 0
	for i := 0; iterator.Scan(); i++ {
		d := iterator.Doc()
		queue <- d

		fmt.Printf("\rdocs: %d, tokens: ---", i)
		// if i == 10 {
		// 	break
		// }
	}

	close(queue)
	for i := 0; i < workers; i++ {
		<-ready
	}

	for i := 0; i < indexWriterCount; i++ {
		if err := indexes[i].Close(); err != nil {
			panic(err)
		}
	}

	// leng := 0
	// for k, v := range tokenStat {
	// 	fmt.Println(k, v)
	// 	leng += len(k)
	// }

	// fmt.Println("unique tokens:", len(tokenStat))
	// fmt.Println("tokens:", tokensCount)
	// fmt.Println("average tokens:", tokensCount/len(tokenStat))
	// fmt.Println("length tokens:", leng)

	// fmt.Printf("lexemizer stat: all - %.3f -> unique %.3f -> words only %.3f -> long words %.3f -> stem %.3f\n",
	// 	float64(lexemizer.StatAll)/float64(lexemizer.StatAll),
	// 	float64(lexemizer.StatUnique)/float64(lexemizer.StatAll),
	// 	float64(lexemizer.StatWordsOnly)/float64(lexemizer.StatAll),
	// 	float64(lexemizer.StatLongWords)/float64(lexemizer.StatAll),
	// 	float64(lexemizer.StatStem)/float64(lexemizer.StatAll))

	// fmt.Println("iteration error:", iterator.Err())
}

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return (int)(h.Sum32())
}

func worker() {
	lexemizer := lexeme.NewParser()

	for d := range queue {
		tokens := lexemizer.Parse(d.String())
		for _, token := range tokens {
			i := hash(token) % indexWriterCount
			if err := indexes[i].AddLexeme(d.DocInfo.ID, token); err != nil {
				fmt.Println("error adding token:", err)
				return
			}
			if i == -1 {
				panic(i)
			}
		}
	}
	ready <- true
}
