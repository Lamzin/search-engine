package numberscompressor

import (
	"math/bits"
)

type EliasGammaCodes struct{}

func (*EliasGammaCodes) Compress(numbers []uint32) (bytes []byte, err error) {
	bitPayload := 0
	for _, number := range numbers {
		leadingZeros := bits.Len32(number) - 1
		requiredAdditionalBytes := (bitPayload+leadingZeros+(leadingZeros+1)+7)/8 - len(bytes)
		bytes = append(bytes, make([]byte, requiredAdditionalBytes)...)
		bitPayload += leadingZeros
		for i := leadingZeros; i >= 0; i-- {
			if number&(1<<(uint)(i)) > 0 {
				bytes[bitPayload/8] |= 1 << (uint)(7-bitPayload%8)
			}
			bitPayload++
		}
	}
	return
}

func (*EliasGammaCodes) Decompress(data []byte) ([]uint32, error) {
	numbers := make([]uint32, 0)
	bits := len(data) * 8
	for i := 0; i < bits; {
		leadingZeros := 0
		for ; i < bits && data[i/8]&(1<<(uint)(7-i%8)) == 0; leadingZeros++ {
			i++
		}
		if i+leadingZeros >= bits {
			break
		}
		var number uint32
		for j := i; j < i+leadingZeros+1; j++ {
			if data[j/8]&(1<<(uint)(7-j%8)) > 0 {
				number |= 1 << (uint)(i+leadingZeros-j)
			}
		}
		numbers = append(numbers, number)
		i += leadingZeros + 1
	}
	return numbers, nil
}
