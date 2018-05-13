package indexwriter

import (
	"github.com/lamzin/search-engine/algos/compressor/numbers"
)

var (
	// variableByteCodes = numberscompressor.VariableByteCodes{}
	variableByteCodes = numberscompressor.BigEndian{}
	bigEndian         = numberscompressor.BigEndian{}
)

type Writer interface {
	AddLexeme(docID int, lexeme string) error
	Close() error
}
