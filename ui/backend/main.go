package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lamzin/search-engine/algos/compressor/text"
	"github.com/lamzin/search-engine/algos/lexeme"
	"github.com/lamzin/search-engine/doc"
	"github.com/lamzin/search-engine/index/reader"
)

var (
	index indexreader.Reader

	articlesPath string
	indexPath    string

	parser *lexeme.Parser = lexeme.NewParser()
)

type searchResult struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	ShortBody string `json:"short_body"`
}

func main() {
	if len(os.Args) != 3 {
		panic("usage: cmd path/to/articles path/to/index")
	}
	articlesPath, indexPath = os.Args[1], os.Args[2]
	fmt.Printf("Articles root: %s\n", articlesPath)
	fmt.Printf("Index root: %s\n", indexPath)

	go initIndex()

	http.HandleFunc("/search", search)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if index == nil {
		fmt.Fprintln(w, "index not ready")
		return
	}

	query, ok := r.URL.Query()["query"]
	if !ok {
		fmt.Fprintf(w, "empty query")
		return
	}
	if len(query) != 1 {
		fmt.Fprintf(w, "invalid query %s", query)
		return
	}

	fmt.Println(r.URL.Query())

	skip, _ := strconv.Atoi(r.URL.Query()["skip"][0])
	limit, _ := strconv.Atoi(r.URL.Query()["limit"][0])

	docs, count, err := findDoc(query[0], skip, limit)
	if err != nil {
		panic(err)
	}
	results := make([]searchResult, len(docs))
	for i, d := range docs {
		results[i].ID = d.ID
		results[i].Title = d.Name
		results[i].ShortBody = d.String()[:200]
	}

	response := map[string]interface{}{
		"results_count": count,
		"results":       results,
	}

	js, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func initIndex() {
	var err error
	index, err = indexreader.NewDBReader(indexPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("index ready")
}

func findDoc(text string, skip int, limit int) ([]*doc.Doc, int, error) {
	lexemes := parser.Parse(text)
	fmt.Println(text, "-->", lexemes)

	lexeme := lexemes[0]

	docIDs, err := index.GetDocIDs(lexeme)
	if err != nil {
		return nil, 0, err
	}
	docCount := len(docIDs)
	docIDs = docIDs[skip : skip+limit]
	fmt.Println("doc ids:", docIDs)

	var docs []*doc.Doc

	for _, docID := range docIDs {
		found := false
		docreader := doc.NewDocCompressedReader(filepath.Join(articlesPath, strconv.Itoa(docID/100)), textcompressor.GzipCompressor{})
		for i := 0; docreader.Scan(); i++ {
			d := docreader.Doc()
			if i == docID%100 {
				d.ID = docID
				docs = append(docs, d)
				found = true
				break
			}
		}
		docreader.Close()

		if err := docreader.Err(); err != nil {
			return nil, 0, err
		}
		if !found {
			return nil, 0, fmt.Errorf("doc not found: %d", docID)
		}
	}
	return docs, docCount, nil
}
