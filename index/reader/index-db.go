package indexreader

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lamzin/search-engine/index/common"
)

const (
	startBuffer = 4
	maxFiles    = 26
)

type IndexDBReader struct {
	indexPath string

	lexemeInfo map[string]*indexcommon.LexemeInfo

	files []*os.File

	lock chan bool
}

func NewDBReader(indexPath string) (*IndexDBReader, error) {
	file, err := os.Open(filepath.Join(indexPath, "lexeme"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}

	content := buf.Bytes()
	// bytes, err := ioutil.ReadFile(filepath.Join(indexPath, "lexeme"))
	// if err != nil {
	// 	return nil, err
	// }

	lexemeInfo := make(map[string]*indexcommon.LexemeInfo)

	for i := 0; i < len(content); i += 4 {
		numbers, err := bigEndian.Decompress(content[i : i+4])
		if err != nil {
			return nil, err
		}
		info, err := indexcommon.NewLexemeInfo(content[i+4 : i+4+numbers[0]])
		if err != nil {
			return nil, err
		}
		lexemeInfo[info.Lexeme] = info
		i += numbers[0]
	}

	files := make([]*os.File, maxFiles)
	for i := 0; i < maxFiles; i++ {
		f, err := os.OpenFile(filepath.Join(indexPath, strconv.Itoa(i)), os.O_RDONLY, 777)
		if err != nil {
			return nil, err
		}
		files[i] = f
	}

	return &IndexDBReader{
		indexPath:  indexPath,
		lexemeInfo: lexemeInfo,
		files:      files,
	}, nil
}

func (r *IndexDBReader) GetDocIDs(lexeme string) ([]int, error) {
	info, ok := r.lexemeInfo[lexeme]
	if !ok {
		return nil, fmt.Errorf("lexeme not found: %s", lexeme)
	}
	docIDs := make([]int, 0)
	for i, position := range info.Positions {
		size := startBuffer << (uint)(i)
		if i+1 == len(info.Positions) {
			size = info.LastLength
		}
		bytes := make([]byte, size)
		if _, err := r.files[i].ReadAt(bytes, (int64)(position)); err != nil {
			return nil, err
		}
		numbers, err := bigEndian.Decompress(bytes)
		if err != nil {
			return nil, err
		}
		docIDs = append(docIDs, numbers...)
	}
	return docIDs, nil
}
