package numberscompressor

import (
	"encoding/binary"
)

type BigEndian struct{}

func (b *BigEndian) Compress(numbers []int) ([]byte, error) {
	var buffer []byte
	bytes := make([]byte, 4)
	for _, number := range numbers {
		binary.BigEndian.PutUint32(bytes, uint32(number))
		buffer = append(buffer, bytes...)
	}
	return buffer, nil
}

func (b *BigEndian) Decompress(data []byte) ([]int, error) {
	var numbers []int
	for i := 0; i < len(data); i += 4 {
		number := int(binary.BigEndian.Uint32(data[i : i+4]))
		numbers = append(numbers, number)
	}
	return numbers, nil
}
