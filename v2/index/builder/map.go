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

	docCount         int
	lexemeCount      int
	lexemeCountQueue chan int
}

func NewIndexBuilderMapper(docPath string, indexPath string) *IndexBuilderMapper {
	return &IndexBuilderMapper{
		docPath:   docPath,
		indexPath: indexPath,

		taskQueue:        make(chan *doc.Doc, mapQueueSize),
		ready:            make(chan struct{}, 1),
		lexemeCountQueue: make(chan int, 1),
	}
}

func (m *IndexBuilderMapper) Run() error {
	go m.lexemeCounter()
	for i := 0; i < mapWorkers; i++ {
		go m.worker()
	}

	iterator := doc.NewDocAllIterator(m.docPath)
	for i := 0; iterator.Scan(); i++ {
		d := iterator.Doc()
		m.taskQueue <- d
	}
	close(m.taskQueue)
	for i := 0; i < mapWorkers; i++ {
		<-m.ready
	}

	close(m.lexemeCountQueue)
	<-m.ready

	fmt.Printf("\nmap finish successfull: docs %d, lexemes: %d\n\n", m.docCount, m.lexemeCount)
	return nil
}

func (m *IndexBuilderMapper) lexemeCounter() {
	for count := range m.lexemeCountQueue {
		m.docCount++
		m.lexemeCount += count
		fmt.Printf("\rdocs: %d, lexemes: %d", m.docCount, m.lexemeCount)
	}
	m.ready <- struct{}{}
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
		m.lexemeCountQueue <- len(wordFrequencies)
	}

	if err := index.Dump(); err != nil {
		panic(err)
	}

	m.ready <- struct{}{}
}
