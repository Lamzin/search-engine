package inverted

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"encoding/binary"

	"github.com/lamzin/search-engine/index/common"
	"github.com/lamzin/search-engine/index/model/doc"
)

type InvertedIndex struct {
	IndexPath string
}

func NewInvertedIndex(indexPath string) *InvertedIndex {
	return &InvertedIndex{
		IndexPath: indexPath,
	}
}

func (i *InvertedIndex) GetDocIDs(token string) ([]int, error) {
	_, filePath := common.FilePath(token)

	lines, err := common.ReadFile(filepath.Join(i.IndexPath, filePath))
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

func (i *InvertedIndex) AddToken(info *doc.DocInfo, token string) error {
	bs := make([]byte, 4)
    binary.LittleEndian.PutUint32(bs, uint32(info.ID))
	
	_, filePath := common.FilePath(token)
	// return common.AppendFile(filepath.Join(i.IndexPath, filePath), fmt.Sprintf("%d", info.ID))
	return common.AppendFile(filepath.Join(i.IndexPath, filePath), string(bs))
}
