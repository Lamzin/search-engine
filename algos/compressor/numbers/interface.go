package numberscompressor

type Compressor interface {
	Compress(numbers []int) ([]byte, error)
	Decompress(data []byte) ([]int, error)
}
