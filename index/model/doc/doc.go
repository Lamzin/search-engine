package doc

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	separator = " ### "
)

type DocInfo struct {
	ID   int
	Name string
	Path string
}

func (doc *DocInfo) String() string {
	return strings.Join(
		[]string{strconv.FormatInt((int64)(doc.ID), 10), doc.Path, doc.Name},
		separator)
}

func DocInfoFromString(s string) (*DocInfo, error) {
	parts := strings.Split(s, separator)
	if len(parts) < 3 {
		return nil, fmt.Errorf("not enough separators", s)
	}

	id, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return nil, err
	}
	return &DocInfo{
		ID:   (int)(id),
		Name: strings.Join(parts[2:], separator),
		Path: parts[1],
	}, nil
}

type Doc struct {
	DocInfo
	Lines []string
}

func (doc *Doc) AddLine(line string) {
	doc.Lines = append(doc.Lines, line)
}

func (doc *Doc) String() string {
	return strings.Join(append([]string{doc.Name, "\n"}, doc.Lines...), "\n")
}
