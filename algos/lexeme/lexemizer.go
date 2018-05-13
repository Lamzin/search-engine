package lexeme

import (
	"regexp"
	"sort"
	"strings"

	"github.com/reiver/go-porterstemmer"
)

var (
	tokenRegEx, errTokenRegEx = regexp.Compile("[a-z0-9-']+")
	wordRegEx, errWordRegEx   = regexp.Compile("^[a-z-']+$")
)

type Parser struct {
	text    string
	lexemes []string

	// stat
	StatAll       int
	StatWordsOnly int
	StatLongWords int
	StatUnique    int
	StatStem      int
}

func NewParser() *Parser {
	return &Parser{}
}

func (l *Parser) Parse(text string) []string {
	l.text = text
	l.all()
	l.StatAll += len(l.lexemes)
	l.unique()
	l.StatUnique += len(l.lexemes)
	l.wordsOnly()
	l.StatWordsOnly += len(l.lexemes)
	l.longWords()
	l.StatLongWords += len(l.lexemes)
	l.stem()
	l.StatStem += len(l.lexemes)
	return l.lexemes
}

func (l *Parser) all() {
	l.lexemes = tokenRegEx.FindAllString(strings.ToLower(l.text), -1)
}

func (l *Parser) wordsOnly() {
	for i := 0; i < len(l.lexemes); i++ {
		if !wordRegEx.MatchString(l.lexemes[i]) {
			l.lexemes = append(l.lexemes[:i], l.lexemes[i+1:]...)
			i--
		}
	}
}

func (l *Parser) longWords() {
	for i := 0; i < len(l.lexemes); i++ {
		if len(l.lexemes[i]) < 3 {
			l.lexemes = append(l.lexemes[:i], l.lexemes[i+1:]...)
			i--
		}
	}
}

func (l *Parser) unique() {
	sort.Strings(l.lexemes)
	for i := 1; i < len(l.lexemes); i++ {
		if l.lexemes[i] == l.lexemes[i-1] {
			l.lexemes = append(l.lexemes[:i], l.lexemes[i+1:]...)
			i--
		}
	}
}

func (l *Parser) stem() {
	for i := 0; i < len(l.lexemes); i++ {
		l.lexemes[i] = porterstemmer.StemString(l.lexemes[i])
	}
	l.unique()
}
