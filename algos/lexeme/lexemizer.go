package lexeme

import (
	"regexp"
	"strings"

	"github.com/reiver/go-porterstemmer"
)

var wordRegEx, errWordRegEx = regexp.Compile("[a-z-']{3,}")

type WordFrequency struct {
	Word      string
	Frequency uint32
}

type Parser struct{}

func (l *Parser) Parse(text string) (frequencies []WordFrequency) {
	words := l.getWords(text)
	wordToFrequency := make(map[string]uint32)
	for _, word := range words {
		wordToFrequency[l.stem(word)]++
	}
	for word, frequency := range wordToFrequency {
		frequencies = append(frequencies, WordFrequency{word, frequency})
	}
	return
}

func (l *Parser) getWords(text string) []string {
	return wordRegEx.FindAllString(strings.ToLower(text), -1)
}

func (l *Parser) stem(lexeme string) string {
	return porterstemmer.StemString(lexeme)
}
