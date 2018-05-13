package lib

import (
	"github.com/lamzin/search-engine/index/inverted"
	"github.com/lamzin/search-engine/index/model/doc"
)

type InvertedIndexWriter struct {
	docManager *doc.DocManager

	lexemizer *Lexemizer

	index *inverted.InvertedIndex
}

// func NewInvertedIndexWriter(articlesPath string, indexPath string) (*InvertedIndexWriter, error) {
// 	docManager, err := doc.NewDocManager(articlesPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	lexemizer, err := NewLexemizer()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &InvertedIndexWriter{
// 		docManager: docManager,
// 		lexemizer:  lexemizer,
// 		index:      inverted.NewInvertedIndex(indexPath),
// 	}, nil
// }

// func (w *InvertedIndexWriter) Run() error {
// 	docInfos, err := w.docManager.GetAllList()
// 	if err != nil {
// 		return err
// 	}

// 	for i, info := range docInfos {
// 		fmt.Printf("\rProgress %d/%d", i, len(docInfos))
// 		doc, err := w.docManager.GetByInfo(info)
// 		if err != nil {
// 			return err
// 		}
// 		tokens := w.lexemizer.Parse(doc.String())
// 		// if len(tokens) < 30 {
// 		// 	fmt.Println(doc.Lines)
// 		// 	fmt.Println(tokens)
// 		// }
// 		for _, token := range tokens {
// 			err := w.index.AddToken(&info, token)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }
