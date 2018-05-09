package doc

import (
	"path/filepath"
	"strconv"

	"github.com/lamzin/search-engine/algos/compressor"
)

const (
	chunkSize = 100
)

type DocManager interface {
	DumpDocument(doc *Doc) error
	Close() error
}

type DocFileManager struct {
	docRoot string

	docWriter DocWriter
	comp      compressor.Compressor
	chunks    int
}

func NewDocFileManager(docRoot string) *DocFileManager {
	return &DocFileManager{
		docRoot: docRoot,
		comp:    &compressor.GzipCompressor{Level: compressor.BestCompression},
	}
}

func (m *DocFileManager) DumpDocument(doc *Doc) error {
	if m.docWriter == nil || (m.docWriter != nil && m.docWriter.Count() >= chunkSize) {
		if err := m.newWriter(); err != nil {
			return err
		}
	}
	return m.docWriter.Write(doc)
}

func (m *DocFileManager) Close() error {
	if m.docWriter != nil {
		return m.docWriter.Close()
	}
	return nil
}

func (m *DocFileManager) newWriter() error {
	if m.docWriter != nil {
		if err := m.docWriter.Close(); err != nil {
			return err
		}
		m.chunks++
	}
	writer, err := NewCompressWriter(filepath.Join(m.docRoot, strconv.FormatInt((int64)(m.chunks), 10)), m.comp)
	m.docWriter = writer
	return err
}
