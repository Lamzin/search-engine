package doc

import (
	"bytes"
	"fmt"
	"os"

	"github.com/lamzin/search-engine/algos/compressor/text"
)

type DocWriter interface {
	Write(doc *Doc) error
	Close() error
	Count() int
}

type DocCompressWriter struct {
	file   *os.File
	comp   textcompressor.Compressor
	buffer bytes.Buffer
	count  int
}

func NewCompressWriter(filePath string, comp textcompressor.Compressor) (*DocCompressWriter, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return &DocCompressWriter{
		file: file,
		comp: comp,
	}, nil
}

func (w *DocCompressWriter) Write(doc *Doc) error {
	w.count++
	if _, err := w.buffer.WriteString(doc.String()); err != nil {
		return err
	}
	_, err := w.buffer.WriteString("\n")
	return err
}

func (w *DocCompressWriter) Close() error {
	defer w.file.Close()

	b := w.buffer.Bytes()
	compressed, err := w.comp.Compress(b)
	if err != nil {
		return err
	}
	_, err = w.file.Write(compressed)

	fmt.Printf("%dK --> %dK, %d persents\n", len(b)/1000, len(compressed)/1000, 100*len(compressed)/len(b))
	return err
}

func (w *DocCompressWriter) Count() int {
	return w.count
}
