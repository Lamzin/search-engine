package lib

import (
	"regexp"
	"strings"
)

type Lexemizer struct {
	regEx *regexp.Regexp
}

func NewLexemizer() (*Lexemizer, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9- ]+")
	if err != nil {
		return nil, err
	}
	return &Lexemizer{
		regEx: reg,
	}, nil
}

func (l *Lexemizer) Parse(text string) []string {
	text = l.regEx.ReplaceAllString(strings.ToLower(text), " ")
	text = strings.Replace(text, "  ", " ", -1)
	words := make([]string, 0)
	for _, word := range strings.Split(text, " ") {
		if len(word) > 0 {
			words = append(words, word)
		}
	}
	return words
}
