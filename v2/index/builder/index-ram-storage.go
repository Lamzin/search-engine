package builder

import (
	"fmt"
	"os"
)

type IndexRAMStorage struct {
	Infos         []LexemeStorageInfo
	invertedInfos map[string]int

	postingsFile    *os.File
	frequenciesFile *os.File

	postingsSize    int64
	frequenciesSize int64
}

func NewIndexRAMStorage(indexPath string) (*IndexRAMStorage, error) {
	postingsFile, postingsSize, err := openFileToRead(indexPath + EXT_POSTINGS)
	if err != nil {
		return nil, err
	}

	frequenciesFile, frequenciesSize, err := openFileToRead(indexPath + EXT_FREQUENCIES)
	if err != nil {
		postingsFile.Close()
		return nil, err
	}

	var infos LexemeStorageInfos
	if err := infos.Load(indexPath); err != nil {
		postingsFile.Close()
		frequenciesFile.Close()
		return nil, err
	}

	invertedInfos := make(map[string]int, 0)
	for i, info := range infos {
		invertedInfos[info.Lexeme] = i
	}

	return &IndexRAMStorage{
		Infos:         infos,
		invertedInfos: invertedInfos,

		postingsFile:    postingsFile,
		frequenciesFile: frequenciesFile,

		postingsSize:    postingsSize,
		frequenciesSize: frequenciesSize,
	}, nil
}

func (index *IndexRAMStorage) GetPostingsAndFrequencies(lexeme string) (*twoArrays, error) {
	postings, err := index.GetPostings(lexeme)
	if err != nil {
		return nil, err
	}
	frequencies, err := index.GetFrequencies(lexeme)
	if err != nil {
		return nil, err
	}
	return &twoArrays{Key: postings, Value: frequencies}, nil
}

func (index *IndexRAMStorage) GetPostings(lexeme string) ([]uint32, error) {
	i, ok := index.invertedInfos[lexeme]
	if !ok {
		return nil, fmt.Errorf("postings do not exist for lexeme: %s", lexeme)
	}
	var start, end = (int64)(index.Infos[i].PostingsStartAt), index.postingsSize
	if i+1 < len(index.Infos) {
		end = (int64)(index.Infos[i+1].PostingsStartAt)
	}
	buffer := make([]byte, end-start)
	if _, err := index.postingsFile.ReadAt(buffer, start); err != nil {
		return nil, err
	}
	deltas, err := bytesCoding.Decompress(buffer)
	if err != nil {
		return nil, err
	}
	return deltaCoding.Encode(deltas), nil
}

func (index *IndexRAMStorage) GetFrequencies(lexeme string) ([]uint32, error) {
	i, ok := index.invertedInfos[lexeme]
	if !ok {
		return nil, fmt.Errorf("frequencies do not exist for lexeme: %s", lexeme)
	}
	var start, end = (int64)(index.Infos[i].FrequenciesStartAt), index.frequenciesSize
	if i+1 < len(index.Infos) {
		end = (int64)(index.Infos[i+1].FrequenciesStartAt)
	}
	buffer := make([]byte, end-start)
	if _, err := index.frequenciesFile.ReadAt(buffer, start); err != nil {
		return nil, err
	}
	return eliasCoding.Decompress(buffer)
}

func (a *IndexRAMStorage) Merge(b *IndexRAMStorage, destinationPath string) error {
	postingsFile, err := os.OpenFile(destinationPath+EXT_POSTINGS, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	frequenciesFile, err := os.OpenFile(destinationPath+EXT_FREQUENCIES, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	var cInfos LexemeStorageInfos
	var cPostingsFileSize, cFrequenciesFileSize uint32

	aiter, biter := 0, 0
	for aiter+biter < len(a.Infos)+len(b.Infos) {
		var cArrays *twoArrays
		var alexeme, blexeme, clexeme string
		if aiter < len(a.Infos) {
			alexeme = a.Infos[aiter].Lexeme
		}
		if biter < len(b.Infos) {
			blexeme = b.Infos[biter].Lexeme
		}
		if alexeme == blexeme {
			clexeme = alexeme
			aArrays, err := a.GetPostingsAndFrequencies(alexeme)
			if err != nil {
				return err
			}
			bArrays, err := b.GetPostingsAndFrequencies(blexeme)
			if err != nil {
				return err
			}
			cArrays = aArrays.Merge(bArrays)
			aiter++
			biter++
		} else if (len(alexeme) > 0 && len(blexeme) == 0) || (len(alexeme) > 0 && alexeme < blexeme) {
			clexeme = alexeme
			var err error
			if cArrays, err = a.GetPostingsAndFrequencies(alexeme); err != nil {
				return err
			}
			aiter++
		} else if (len(blexeme) > 0 && len(alexeme) == 0) || (len(blexeme) > 0 && blexeme < alexeme) {
			clexeme = blexeme
			var err error
			if cArrays, err = b.GetPostingsAndFrequencies(blexeme); err != nil {
				return err
			}
			biter++
		} else {
			panic("merge error")
		}

		cInfos = append(cInfos, LexemeStorageInfo{
			Lexeme:             clexeme,
			PostingsStartAt:    cPostingsFileSize,
			FrequenciesStartAt: cFrequenciesFileSize,
		})

		{
			delta := deltaCoding.Decode(cArrays.Key)
			bytes, _ := bytesCoding.Compress(delta)
			cPostingsFileSize += (uint32)(len(bytes))
			if _, err = postingsFile.Write(bytes); err != nil {
				return err
			}
		}

		{
			bytes, _ := eliasCoding.Compress(cArrays.Value)
			cFrequenciesFileSize += (uint32)(len(bytes))
			if _, err = frequenciesFile.Write(bytes); err != nil {
				return err
			}
		}
	}

	return cInfos.Dump(destinationPath)
}
