package builder

import (
	"fmt"

	"github.com/lamzin/search-engine/algos/lexeme"
	"github.com/lamzin/search-engine/doc"
)

const (
	mapWorkers   = 4
	mapQueueSize = 100
)

var lexemizer lexeme.Parser

type IndexBuilderMapper struct {
	docPath   string
	indexPath string

	taskQueue chan *doc.Doc
	ready     chan struct{}
}

func NewIndexBuilderMapper(docPath string, indexPath string) *IndexBuilderMapper {
	return &IndexBuilderMapper{
		docPath:   docPath,
		indexPath: indexPath,

		taskQueue: make(chan *doc.Doc, mapQueueSize),
		ready:     make(chan struct{}, 1),
	}
}

func (m *IndexBuilderMapper) Run() error {
	for i := 0; i < mapWorkers; i++ {
		go m.worker()
	}

	iterator := doc.NewDocAllIterator(m.docPath)
	for i := 0; iterator.Scan(); i++ {
		d := iterator.Doc()
		m.taskQueue <- d
		fmt.Printf("\rdocs: %d", i)
	}
	close(m.taskQueue)

	for i := 0; i < mapWorkers; i++ {
		<-m.ready
	}

	return nil
}

func (m *IndexBuilderMapper) worker() {
	index := NewIndexRAM(m.indexPath)

	for d := range m.taskQueue {
		wordFrequencies := lexemizer.Parse(d.String())
		for _, wordFrequency := range wordFrequencies {
			if index.CanAddLexeme() {
				if err := index.AddLexeme(wordFrequency.Word, d.ID, wordFrequency.Frequency); err != nil {
					panic(err)
				}
			} else {
				fmt.Println("start index dump...")
				if err := index.Dump(); err != nil {
					panic(err)
				}
				fmt.Println("finish index dump")
				index = NewIndexRAM(m.indexPath)
			}
		}
	}

	if err := index.Dump(); err != nil {
		panic(err)
	}

	m.ready <- struct{}{}
}
