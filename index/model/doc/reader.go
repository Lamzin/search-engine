package doc

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lamzin/search-engine/algos/compressor"
)

type DocReader interface {
	Scan() bool
	Doc() *Doc
	Err() error
	Close()
}

type DocTxtReader struct {
	file     *os.File
	scanner  *bufio.Scanner
	finished bool

	line string

	err error
}

func NewDocTxtReader(docPath string) *DocTxtReader {
	file, err := os.Open(docPath)
	return &DocTxtReader{
		file:    file,
		scanner: bufio.NewScanner(file),
		err:     err,
	}
}

func (r *DocTxtReader) Scan() bool {
	return r.err == nil && !r.finished
}

func (r *DocTxtReader) Doc() *Doc {
	if len(r.line) == 0 && r.scanner.Scan() {
		r.line = r.scanner.Text()
	}

	doc := Doc{
		DocInfo: DocInfo{
			Name: r.line,
		},
		Lines: []string{r.line},
	}

	for r.scanner.Scan() {
		r.line = r.scanner.Text()
		if strings.HasPrefix(r.line, "= ") {
			return &doc
		}
		doc.AddLine(r.line)
	}
	r.finished = true
	return &doc
}

func (r *DocTxtReader) Err() error {
	if r.err != nil {
		return r.err
	}
	return r.scanner.Err()
}

func (r *DocTxtReader) Close() {
	r.file.Close()
}

type DocCompressedReader struct {
	filePath string

	docs     []*Doc
	docIndex int

	comp compressor.Compressor

	err error
}

func NewDocCompressedReader(docPath string, comp compressor.Compressor) *DocCompressedReader {
	return &DocCompressedReader{
		filePath: docPath,
		comp:     comp,
	}
}

func (r *DocCompressedReader) Scan() bool {
	if r.docIndex == 0 {
		var b []byte
		b, r.err = ioutil.ReadFile(r.filePath)
		if r.err != nil {
			return false
		}
		b, r.err = r.comp.Decompress(b)
		if r.err != nil {
			return false
		}

		lines := strings.Split(string(b), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "= ") {
				r.docs = append(r.docs, &Doc{DocInfo: DocInfo{Name: line}})
			}
			if len(r.docs) == 0 {
				r.err = fmt.Errorf("empty file or wrong format")
				return false
			}
			r.docs[len(r.docs)-1].AddLine(line)
		}
	}
	return r.err == nil && r.docIndex < len(r.docs)
}

func (r *DocCompressedReader) Doc() (d *Doc) {
	if r.docIndex < len(r.docs) {
		d = r.docs[r.docIndex]
	}
	r.docIndex++
	return
}

func (r *DocCompressedReader) Err() error {
	return r.err
}

func (r *DocCompressedReader) Close() {}
