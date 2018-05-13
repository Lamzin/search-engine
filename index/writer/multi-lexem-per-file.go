package indexwriter

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/lamzin/search-engine/utils"
)

const (
	multiFileSizeLimit = 4096
)

type MultiLexemePerFile struct {
	indexPath string

	lexemeToFile map[string]int
	fileToLexeme [][]string
	fileSizes    []int
}

func NewMultiLexemePerFile(indexPath string) *MultiLexemePerFile {
	return &MultiLexemePerFile{
		indexPath:    indexPath,
		lexemeToFile: make(map[string]int, 0),
	}
}

func (w *MultiLexemePerFile) AddLexeme(docID int, lexeme string) error {
	file := w.getFileID(lexeme)
	bytes, err := utils.ReadFileBytes(filepath.Join(w.indexPath, strconv.Itoa(file)))
	if err != nil {
		return err
	}
	// fmt.Println("bytes", bytes)
	postings, err := w.bytesToPostings(bytes)
	if err != nil {
		return err
	}
	// fmt.Println("postings:", postings)
	// fmt.Println("file:", file)
	// fmt.Println("lexeme:", lexeme)
	// fmt.Println("w.fileToLexeme[file]:", w.fileToLexeme[file])

	for i, l := range w.fileToLexeme[file] {
		if l == lexeme {
			if i == len(postings) {
				postings = append(postings, []int{docID})
			} else {
				postings[i] = append(postings[i], docID)
			}
			break
		}
	}

	bytes, err = w.postingsToBytes(postings)
	if err != nil {
		return err
	}
	w.fileSizes[file] = len(bytes)
	return utils.WriteFile(filepath.Join(w.indexPath, strconv.Itoa(file)), bytes)
}

func (w *MultiLexemePerFile) Close() error {
	return nil
}

func (w *MultiLexemePerFile) getFileID(lexeme string) int {
	if file, ok := w.lexemeToFile[lexeme]; ok {
		return file
	}
	for file, size := range w.fileSizes {
		if size < multiFileSizeLimit {
			w.fileToLexeme[file] = append(w.fileToLexeme[file], lexeme)
			w.lexemeToFile[lexeme] = file
			return file
		}
	}
	file := len(w.fileToLexeme)
	w.fileToLexeme = append(w.fileToLexeme, []string{lexeme})
	w.lexemeToFile[lexeme] = file
	w.fileSizes = append(w.fileSizes, 0)
	return file
}

func (w *MultiLexemePerFile) bytesToPostings(bytes []byte) ([][]int, error) {
	var postings [][]int
	index := 0
	for index < len(bytes) {
		numbers, err := bigEndian.Decompress(bytes[index : index+4])
		if err != nil {
			return nil, err
		}
		dataLength := numbers[0]
		if index+4+dataLength > len(bytes) {
			fmt.Println(index, dataLength, index+4+dataLength, bytes)
		}
		numbers, err = variableByteCodes.Decompress(bytes[index+4 : index+4+dataLength])
		if err != nil {
			return nil, err
		}
		postings = append(postings, numbers)
		index += 4 + dataLength
	}
	return postings, nil
}

func (w *MultiLexemePerFile) postingsToBytes(postings [][]int) ([]byte, error) {
	var bytes []byte
	for _, ids := range postings {
		data, err := variableByteCodes.Compress(ids)
		if err != nil {
			return nil, err
		}
		dataLength, err := bigEndian.Compress([]int{len(data)})
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, dataLength...)
		bytes = append(bytes, data...)
	}
	return bytes, nil
}
