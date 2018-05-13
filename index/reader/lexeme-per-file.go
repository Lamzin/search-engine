package indexreader

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lamzin/search-engine/utils"
)

type LexemePerFile struct {
	indexPath string
}

func NewLexemPerFile(indexPath string) *LexemePerFile {
	return &LexemePerFile{indexPath: indexPath}
}

func (r *LexemePerFile) GetDocIDs(lexeme string) ([]int, error) {
	_, filePath := utils.FilePath(lexeme)

	lines, err := utils.ReadFile(filepath.Join(r.indexPath, filePath))
	if err != nil {
		return nil, err
	}
	if len(lines) != 1 {
		return nil, fmt.Errorf("does not contain only one line")
	}

	numbers := strings.Split(strings.TrimSpace(lines[0]), " ")
	var arr []int
	for _, number := range numbers {
		n, err := strconv.ParseInt(number, 10, 32)
		if err != nil {
			return nil, err
		}
		arr = append(arr, (int)(n))
	}
	return arr, nil
}
