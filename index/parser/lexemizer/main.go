package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/lamzin/search-engine/index/model"
)

const (
	lexemizerWorkers = 32
)

func main() {
	if len(os.Args) != 2 {
		panic("usage: cmd path/to/engine/root")
	}
	var engineRoot = os.Args[1]
	fmt.Printf("Engine root: %s\n", engineRoot)

	indexData := model.NewDocProvider(engineRoot)

	lexemizer := &Lexemizer{
		DataRoot:    "",
		DocProvider: indexData,
	}
	if err := lexemizer.Run(); err != nil {
		fmt.Println(err.Error())
	}
}

// Lexem type
type Lexem string

type Lexemizer struct {
	DataRoot    string
	DocProvider model.DocProviderI

	tasksQueue   chan string
	workersReady chan struct{}

	regEx *regexp.Regexp
}

func (l *Lexemizer) Run() error {
	reg, err := regexp.Compile("[^a-zA-Z0-9- ]+")
	if err != nil {
		return err
	}
	l.regEx = reg

	docNames, err := l.DocProvider.GetAllNames()
	if err != nil {
		return err
	}

	l.tasksQueue = make(chan string, len(docNames))
	l.workersReady = make(chan struct{}, 1)
	for _, docName := range docNames {
		l.tasksQueue <- docName
	}
	close(l.tasksQueue)

	for i := 0; i < lexemizerWorkers; i++ {
		go l.worker()
	}

	for i := 0; i < lexemizerWorkers; i++ {
		<-l.workersReady
	}
	return nil
}

func (l *Lexemizer) worker() {
	count := 0
	for docName := range l.tasksQueue {
		lines, err := l.DocProvider.GetByName(docName)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		for _, line := range lines {
			// fmt.Println(line)
			line = strings.Replace(line, "  ", " ", -1)
			line = l.regEx.ReplaceAllString(strings.ToLower(line), "")
			strings.Split(line, " ")
			// words := strings.Split(line, " ")
			// fmt.Println(words)
		}
		count++
		fmt.Printf("count: %d\n", count)
	}
	l.workersReady <- struct{}{}
}
