package numberscompressor

import (
	"fmt"
)

type DeltaCoding struct{}

func (DeltaCoding) Decode(numbers []uint32) []uint32 {
	if numbers == nil {
		return nil
	}
	if len(numbers) == 0 {
		return make([]uint32, 0)
	}
	deltas := make([]uint32, len(numbers))
	deltas[0] = numbers[0]
	for i := 1; i < len(numbers); i++ {
		if numbers[i-1] >= numbers[i] {
			panic(fmt.Sprintf("invalid input: %x", numbers))
		}
		deltas[i] = numbers[i] - numbers[i-1]
	}
	return deltas
}

func (DeltaCoding) Encode(deltas []uint32) []uint32 {
	if deltas == nil {
		return nil
	}
	if len(deltas) == 0 {
		return make([]uint32, 0)
	}
	numbers := make([]uint32, len(deltas))
	numbers[0] = deltas[0]
	for i := 1; i < len(deltas); i++ {
		numbers[i] = numbers[i - 1] + deltas[i]
	}
	return numbers
}
