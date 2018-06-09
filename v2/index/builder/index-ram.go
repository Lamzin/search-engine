package builder

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/lamzin/search-engine/algos/compressor/numbers"
)

const RAM_LIMIT = 10 * 1024 * 1024

// const RAM_LIMIT = 1024

var (
	bigEndian   = numberscompressor.BigEndian{}
	bytesCoding = numberscompressor.VariableByteCodes{}
	deltaCoding = numberscompressor.DeltaCoding{}
	eliasCoding = numberscompressor.EliasGammaCodes{}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type LexemeInfo struct {
	Lexeme      string
	Postings    []uint32
	Frequencies []uint32
}

type lexemeInfoSorter []*LexemeInfo

func (arr lexemeInfoSorter) Len() int {
	return len(arr)
}

func (arr lexemeInfoSorter) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (arr lexemeInfoSorter) Less(i, j int) bool {
	return arr[i].Lexeme < arr[j].Lexeme
}

type IndexRAM struct {
	lexemeInfo map[string]*LexemeInfo
	records    uint32
	workdir    string

	indexFileName string
	infos         []*LexemeInfo
	storageInfos  []LexemeStorageInfo
}

func NewIndexRAM(workdir string) *IndexRAM {
	return &IndexRAM{
		lexemeInfo:    make(map[string]*LexemeInfo, 0),
		records:       0,
		workdir:       workdir,
		indexFileName: strconv.Itoa(rand.Intn(1 << 30)),
	}
}

func (index *IndexRAM) CanAddLexeme() bool {
	return index.records < RAM_LIMIT
}

func (index *IndexRAM) AddLexeme(lexeme string, docID uint32, frequency uint32) error {
	if index.records > RAM_LIMIT {
		return fmt.Errorf("index is full: size = %d", index.records)
	}
	if info, ok := index.lexemeInfo[lexeme]; !ok {
		index.lexemeInfo[lexeme] = &LexemeInfo{
			Lexeme:      lexeme,
			Postings:    []uint32{docID},
			Frequencies: []uint32{frequency},
		}
	} else {
		info.Postings = append(info.Postings, docID)
		info.Frequencies = append(info.Frequencies, frequency)
	}
	index.records++
	return nil
}

func (index *IndexRAM) Dump() error {
	index.dumpPreparation()
	if err := index.dumpPostings(); err != nil {
		return err
	}
	if err := index.dumpFrequencies(); err != nil {
		return err
	}
	if err := index.dumpInfo(); err != nil {
		return err
	}
	return nil
}

func (index *IndexRAM) dumpPreparation() {
	for _, info := range index.lexemeInfo {
		index.infos = append(index.infos, info)
	}
	sort.Sort(lexemeInfoSorter(index.infos))
	index.storageInfos = make([]LexemeStorageInfo, len(index.infos))
	for i, info := range index.infos {
		index.storageInfos[i] = LexemeStorageInfo{Lexeme: info.Lexeme}
	}
}

func (index *IndexRAM) dumpPostings() error {
	var data []byte
	for i, info := range index.infos {
		sort.Sort(MetaArrays{Postings: info.Postings, Frequencies: info.Frequencies})
		delta := deltaCoding.Decode(info.Postings)
		bytes, _ := bytesCoding.Compress(delta)
		index.storageInfos[i].PostingsStartAt = uint32(len(data))
		data = append(data, bytes...)
	}
	return writeBytesToFile(filepath.Join(index.workdir, index.indexFileName+EXT_POSTINGS), data)
}

func (index *IndexRAM) dumpFrequencies() error {
	var data []byte
	for i, info := range index.infos {
		bytes, _ := eliasCoding.Compress(info.Frequencies)
		index.storageInfos[i].FrequenciesStartAt = uint32(len(data))
		data = append(data, bytes...)
	}
	return writeBytesToFile(filepath.Join(index.workdir, index.indexFileName+EXT_FREQUENCIES), data)
}

func (index *IndexRAM) dumpInfo() error {
	infos := LexemeStorageInfos(index.storageInfos)
	return infos.Dump(filepath.Join(index.workdir, index.indexFileName))
}
