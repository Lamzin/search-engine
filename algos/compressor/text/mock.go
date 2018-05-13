package textcompressor

type MockCompressor struct {
}

func (MockCompressor) Compress(data []byte) ([]byte, error) {
	return data, nil
}

func (MockCompressor) Decompress(data []byte) ([]byte, error) {
	return data, nil
}
