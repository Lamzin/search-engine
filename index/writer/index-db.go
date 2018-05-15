package indexwriter

import (
	"os"
	"path/filepath"
	"strconv"
)

const (
	startBuffer = 32
)

type LexemeInfo struct {
	Positions  []int
	LastLength int
}

type IndexDBWriter struct {
	indexPath string

	lexemeInfo map[string]*LexemeInfo

	fileSizes []int
	files     []*os.File
}

func NewIndexDBWriter(indexPath string) *IndexDBWriter {
	return &IndexDBWriter{
		indexPath:  indexPath,
		lexemeInfo: make(map[string]*LexemeInfo, 0),
		fileSizes:  []int{0},
	}
}

func (w *IndexDBWriter) findFileAndPosition(lexeme string) (file int, info *LexemeInfo) {
	info, ok := w.lexemeInfo[lexeme]
	if !ok {
		info = &LexemeInfo{
			Positions:  []int{w.fileSizes[0]},
			LastLength: 0,
		}
		w.lexemeInfo[lexeme] = info
	}

	file = len(info.Positions) - 1
	bufferLength := startBuffer << (uint)(file)
	if info.LastLength == bufferLength {
		file++
		info.LastLength = 0
		for file >= len(w.fileSizes) {
			w.fileSizes = append(w.fileSizes, 0)
		}
		info.Positions = append(info.Positions, w.fileSizes[file])
	} else if info.LastLength > bufferLength {
		panic("length more that allowed buffer size")
	}
	return file, info
}

func (w *IndexDBWriter) writeDocID(docID int, fileInt int, position int) error {
	if fileInt == len(w.files) || w.files[fileInt] == nil {
		filePath := filepath.Join(w.indexPath, strconv.Itoa(fileInt))
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}
		w.files = append(w.files, file)
	}
	file := w.files[fileInt]

	bytes, err := bigEndian.Compress([]int{docID})
	if err != nil {
		return err
	}

	_, err = file.WriteAt(bytes, int64(position))
	return err
}

func (w *IndexDBWriter) AddLexeme(docID int, lexeme string) error {
	file, info := w.findFileAndPosition(lexeme)

	err := w.writeDocID(docID, file, info.Positions[len(info.Positions)-1])
	if err != nil {
		return err
	}
	info.LastLength += 4
	if w.fileSizes[file] == info.Positions[len(info.Positions)-1] {
		w.fileSizes[file] += 4
	}
	return nil
}

func (w *IndexDBWriter) Close() error {
	for _, file := range w.files {
		file.Close()
	}
	return nil
}
