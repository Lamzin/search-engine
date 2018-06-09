package builder

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	mergeWorkers = 4
)

type mergeTask struct {
	First  string
	Second string
}

type IndexBuilderReducer struct {
	indexPath string

	mergeTasks       chan mergeTask
	mergeTaskResults chan string

	workersReady chan struct{}
}

func NewIndexBuilderReducer(indexPath string) *IndexBuilderReducer {
	return &IndexBuilderReducer{
		indexPath:        indexPath,
		mergeTasks:       make(chan mergeTask, mergeWorkers),
		mergeTaskResults: nil,
		workersReady:     make(chan struct{}, mergeWorkers),
	}
}

func (r *IndexBuilderReducer) Run() error {
	indexFiles, err := r.getIndexList(r.indexPath)
	if err != nil {
		return err
	}
	if len(indexFiles) == 0 {
		return fmt.Errorf("no index files")
	}

	r.mergeTaskResults = make(chan string, len(indexFiles))

	for i := 0; i < mergeWorkers; i++ {
		go r.mergeWorker()
	}

	for _, indexFile := range indexFiles {
		r.mergeTaskResults <- indexFile
	}

	for i := 1; i < len(indexFiles); i++ {
		fmt.Printf("\rreduce tasks: %d   ", len(r.mergeTaskResults))
		first := <-r.mergeTaskResults
		second := <-r.mergeTaskResults
		r.mergeTasks <- mergeTask{First: first, Second: second}
	}
	close(r.mergeTasks)

	for i := 0; i < mergeWorkers; i++ {
		<-r.workersReady
	}

	return nil
}

func (r IndexBuilderReducer) mergeWorker() {
	for task := range r.mergeTasks {
		aIndexName, bIndexName := task.First, task.Second
		cIndexName := strconv.Itoa(rand.Intn(1 << 30))
		aIndex, err := NewIndexRAMStorage(filepath.Join(r.indexPath, aIndexName))
		if err != nil {
			panic(err)
		}
		bIndex, err := NewIndexRAMStorage(filepath.Join(r.indexPath, bIndexName))
		if err != nil {
			panic(err)
		}
		if err = aIndex.Merge(bIndex, filepath.Join(r.indexPath, cIndexName)); err != nil {
			panic(err)
		}

		if err = r.removeIndexFiles(aIndexName); err != nil {
			panic(err)
		}
		if err = r.removeIndexFiles(bIndexName); err != nil {
			panic(err)
		}
		r.mergeTaskResults <- cIndexName
	}
	r.workersReady <- struct{}{}
}

func (r IndexBuilderReducer) removeIndexFiles(indexName string) error {
	if err := os.Remove(filepath.Join(r.indexPath, indexName) + EXT_INFO); err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(r.indexPath, indexName) + EXT_POSTINGS); err != nil {
		return err
	}
	return os.Remove(filepath.Join(r.indexPath, indexName) + EXT_FREQUENCIES)
}

func (r *IndexBuilderReducer) getIndexList(dirPath string) ([]string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	uniqueFileNames := make(map[string]struct{}, 0)
	for _, file := range files {
		uniqueFileNames[strings.Split(file.Name(), ".")[0]] = struct{}{}
	}
	var fileNames []string
	for fileName := range uniqueFileNames {
		fileNames = append(fileNames, fileName)
	}
	return fileNames, nil
}
