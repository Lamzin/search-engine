package lib

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

type Lexemizer struct {
	text  string
	words []string

	// stat
	StatAll       int
	StatWordsOnly int
	StatLongWords int
	StatUnique    int
	StatStem      int
}

func NewLexemizer() *Lexemizer {
	return &Lexemizer{}
}

func (l *Lexemizer) Parse(text string) []string {
	l.text = text
	l.all()
	l.StatAll += len(l.words)
	l.unique()
	l.StatUnique += len(l.words)
	l.wordsOnly()
	l.StatWordsOnly += len(l.words)
	l.longWords()
	l.StatLongWords += len(l.words)
	l.stem()
	l.StatStem += len(l.words)
	return l.words
}

func (l *Lexemizer) all() {
	l.words = tokenRegEx.FindAllString(strings.ToLower(l.text), -1)
}

func (l *Lexemizer) wordsOnly() {
	for i := 0; i < len(l.words); i++ {
		if !wordRegEx.MatchString(l.words[i]) {
			l.words = append(l.words[:i], l.words[i+1:]...)
			i--
		}
	}
}

func (l *Lexemizer) longWords() {
	for i := 0; i < len(l.words); i++ {
		if len(l.words[i]) < 3 {
			l.words = append(l.words[:i], l.words[i+1:]...)
			i--
		}
	}
}

func (l *Lexemizer) unique() {
	sort.Strings(l.words)
	for i := 1; i < len(l.words); i++ {
		if l.words[i] == l.words[i-1] {
			l.words = append(l.words[:i], l.words[i+1:]...)
			i--
		}
	}
}

func (l *Lexemizer) stem() {
	for i := 0; i < len(l.words); i++ {
		l.words[i] = porterstemmer.StemString(l.words[i])
	}
	l.unique()
}
