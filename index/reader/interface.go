package indexreader

import (
	"github.com/lamzin/search-engine/algos/compressor/numbers"
)

var (
	bigEndian = numberscompressor.BigEndian{}
)

type Reader interface {
	GetDocIDs(lexeme string) ([]int, error)
}
