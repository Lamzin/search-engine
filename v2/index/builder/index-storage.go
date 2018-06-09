package builder

import (
	"fmt"
	"io/ioutil"

	"github.com/lamzin/search-engine/algos/compressor/text"
)

const (
	EXT_INFO        = ".info"
	EXT_POSTINGS    = ".postings"
	EXT_FREQUENCIES = ".frequencies"
)

var gzip textcompressor.GzipCompressor = textcompressor.GzipCompressor{
	Level: textcompressor.BestCompression,
}

type LexemeStorageInfo struct {
	Lexeme             string
	PostingsStartAt    uint32
	FrequenciesStartAt uint32
}

func (info *LexemeStorageInfo) Decode() []byte {
	data, _ := bigEndian.Compress([]uint32{info.PostingsStartAt, info.FrequenciesStartAt})
	data = append(data, []byte(info.Lexeme)...)
	return data
}

func (info *LexemeStorageInfo) Encode(data []byte) error {
	if len(data) < 8 {
		return fmt.Errorf("too short data: %x", data)
	}
	numbers, _ := bigEndian.Decompress(data[:8])
	info.PostingsStartAt = numbers[0]
	info.FrequenciesStartAt = numbers[1]
	info.Lexeme = string(data[8:])
	return nil
}

type LexemeStorageInfos []LexemeStorageInfo

func (infos *LexemeStorageInfos) Dump(fileName string) error {
	var data []byte
	for _, info := range *infos {
		infoBytes := info.Decode()
		lenInfoBytes, _ := bigEndian.Compress([]uint32{uint32(len(infoBytes))})
		data = append(data, lenInfoBytes...)
		data = append(data, infoBytes...)
	}

	gzipedData, err := gzip.Compress(data)
	if err != nil {
		return err
	}
	return writeBytesToFile(fileName+EXT_INFO, gzipedData)
}

func (infos *LexemeStorageInfos) Load(fileName string) error {
	gzipedData, err := ioutil.ReadFile(fileName + EXT_INFO)
	if err != nil {
		return err
	}

	bytes, err := gzip.Decompress(gzipedData)
	if err != nil {
		return err
	}

	for offset := 0; offset < len(bytes); offset += 4 {
		numbers, _ := bigEndian.Decompress(bytes[offset : offset+4])
		length := (int)(numbers[0])
		var info LexemeStorageInfo
		if err := info.Encode(bytes[offset+4 : offset+4+length]); err != nil {
			return err
		}
		*infos = append(*infos, info)
		offset += length
	}
	return nil
}
