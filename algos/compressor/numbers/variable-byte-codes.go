package numberscompressor

type VariableByteCodes struct{}

func (v *VariableByteCodes) Compress(numbers []int) ([]byte, error) {
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

func (v *VariableByteCodes) Decompress(data []byte) ([]int, error) {
	var numbers []int
	var number int
	for _, b := range data {
		number = number*128 + int(b)%128
		if int(b)&128 == 0 {
			numbers = append(numbers, number)
			number = 0
		}
	}
	return numbers, nil
}
