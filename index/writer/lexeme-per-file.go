package indexwriter

import (
	"path/filepath"

	"github.com/lamzin/search-engine/utils"
)

type LexemePerFile struct {
	indexPath string
}

func NewLexemePerFile(indexPath string) *LexemePerFile {
	return &LexemePerFile{
		indexPath: indexPath,
	}
}

func (w *LexemePerFile) AddLexeme(docID int, lexeme string) error {
	bytes, err := variableByteCodes.Compress([]int{docID})
	if err != nil {
		return err
	}
	_, filePath := utils.FilePath(lexeme)
	return utils.AppendFile(filepath.Join(w.indexPath, filePath), bytes)
}

func (w *LexemePerFile) Close() error {
	return nil
}
