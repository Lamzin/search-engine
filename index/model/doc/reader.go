package doc

import (
	"bufio"
	"os"
	"strings"
)

type DocReader interface {
	Scan() bool
	Doc() *Doc
	Err() error
	Close()
}

type DocTxtReader struct {
	file    *os.File
	scanner *bufio.Scanner

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
	return r.err == nil && r.scanner.Scan()
}

func (r *DocTxtReader) Doc() *Doc {
	if len(r.line) > 0 && r.scanner.Scan() {
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
