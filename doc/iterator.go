package doc

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lamzin/search-engine/algos/compressor/text"
)

type DocAllIterator struct {
	docRoot        string
	docReader      DocReader
	docNumber      uint32
	docReaderIndex uint32

	files []os.FileInfo

	err error
}

func NewDocAllIterator(docRoot string) *DocAllIterator {
	files, err := ioutil.ReadDir(docRoot)

	if err == nil {
		for i := 0; i < len(files); i++ {
			if files[i].IsDir() {
				files = append(files[:i], files[i+1:]...)
				i--
			}
		}
	}
	return &DocAllIterator{
		docRoot: docRoot,
		files:   files,
		err:     err,
	}
}

func (r *DocAllIterator) Scan() bool {
	if r.err != nil {
		return false
	}
	if r.docReader != nil {
		if r.docReader.Scan() {
			return true
		}
		if r.err = r.docReader.Err(); r.err != nil {
			return false
		}
	}
	if len(r.files) > 0 {
		r.docReader = NewDocCompressedReader(filepath.Join(r.docRoot, r.files[0].Name()), textcompressor.GzipCompressor{})
		docNumber, err := strconv.Atoi(r.files[0].Name())
		r.docNumber, r.err = uint32(docNumber), err
			r.docReaderIndex = 0
		r.files = r.files[1:]
		return r.Scan()
	}
	return false
}

func (r *DocAllIterator) Doc() *Doc {
	d := r.docReader.Doc()
	d.ID = uint32(r.docNumber*chunkSize + r.docReaderIndex)
	r.docReaderIndex++
	return d
}

func (r *DocAllIterator) Err() error {
	return r.err
}

func (r *DocAllIterator) Close() {
	if r.docReader != nil {
		r.docReader.Close()
	}
}
