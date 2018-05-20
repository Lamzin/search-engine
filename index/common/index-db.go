package indexcommon

import (
	"github.com/lamzin/search-engine/algos/compressor/numbers"
)

var (
	bigEndian = numberscompressor.BigEndian{}
)

type LexemeInfo struct {
	Lexeme     string
	Positions  []int
	LastLength int
}

func (info *LexemeInfo) Bytes() []byte {
	var bytes []byte
	lenBytes, _ := bigEndian.Compress([]int{len(info.Lexeme)})
	bytes = append(bytes, lenBytes...)
	bytes = append(bytes, []byte(info.Lexeme)...)
	positionsBytes, _ := bigEndian.Compress(info.Positions)
	bytes = append(bytes, positionsBytes...)
	lastLengthBytes, _ := bigEndian.Compress([]int{info.LastLength})
	bytes = append(bytes, lastLengthBytes...)
	return bytes
}

func NewLexemeInfo(bytes []byte) (*LexemeInfo, error) {
	numbers, err := bigEndian.Decompress(bytes[:4])
	if err != nil {
		return nil, err
	}
	lexemeLength := numbers[0]
	numbers, err = bigEndian.Decompress(bytes[4+lexemeLength:])
	if err != nil {
		return nil, err
	}
	return &LexemeInfo{
		Lexeme:     string(bytes[4 : 4+lexemeLength]),
		Positions:  numbers[:len(numbers)-1],
		LastLength: numbers[len(numbers)-1],
	}, nil
}
