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

type IndexBuilderReducer struct {
	indexPath string
}

func NewIndexBuilderReducer(indexPath string) *IndexBuilderReducer {
	return &IndexBuilderReducer{
		indexPath: indexPath,
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

	for len(indexFiles) > 1 {
		aIndexName, bIndexName := indexFiles[0], indexFiles[1]
		cIndexName := strconv.Itoa(rand.Intn(1 << 30))
		indexFiles = append(indexFiles[2:], cIndexName)
		aIndex, err := NewIndexRAMStorage(filepath.Join(r.indexPath, aIndexName))
		if err != nil {
			return err
		}
		bIndex, err := NewIndexRAMStorage(filepath.Join(r.indexPath, bIndexName))
		if err != nil {
			return err
		}
		if err = aIndex.Merge(bIndex, filepath.Join(r.indexPath, cIndexName)); err != nil {
			return err
		}

		if err = r.removeIndexFiles(aIndexName); err != nil {
			return err
		}
		if err = r.removeIndexFiles(bIndexName); err != nil {
			return err
		}
	}
	return nil
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
	for fileName, _ := range uniqueFileNames {
		fileNames = append(fileNames, fileName)
	}
	return fileNames, nil
}
