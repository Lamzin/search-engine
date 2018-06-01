package numberscompressor

type VariableByteCodes struct{}

func (v *VariableByteCodes) Compress(numbers []uint32) ([]byte, error) {
	var buffer []byte
	for _, number := range numbers {
		var bytes []byte
		for number > 127 {
			bytes = append(bytes, byte(number%128))
			number /= 128
		}
		bytes = append(bytes, byte(number))

		for i := len(bytes) - 1; i >= 1; i-- {
			buffer = append(buffer, bytes[i]|128)
		}
		buffer = append(buffer, bytes[0])
	}
	return buffer, nil
}

func (v *VariableByteCodes) Decompress(data []byte) ([]uint32, error) {
	var numbers []uint32
	var number uint32
	for _, b := range data {
		number = number*128 + uint32(b)%128
		if uint32(b)&128 == 0 {
			numbers = append(numbers, number)
			number = 0
		}
	}
	return numbers, nil
}
