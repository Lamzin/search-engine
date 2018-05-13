package indexreader

type Reader interface {
	GetDocIDs(lexeme string) ([]int, error)
}
