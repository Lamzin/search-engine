package numberscompressor

// import (
// 	"math/bits"

// 	"github.com/Workiva/go-datastructures/bitarray"
// )

// type EliasGammaCodes struct{}

// func (*EliasGammaCodes) Compress(numbers []int) ([]byte, error) {
// 	ba := bitarray.NewSparseBitArray()
// 	baLen := 0

// 	var len int
// 	for _, number := range numbers {
// 		len = bits.Len32(uint32(number))
// 		baLen += len
// 		for i := uint(len); i > 0; i-- {
// 			if number&(1<<i) > 0 {
// 				ba.SetBit(uint64(baLen))
// 			}
// 			baLen++
// 		}
// 	}
// 	return nil, nil
// }

// func (*EliasGammaCodes) Decompress(data []byte) ([]int, error) {
// 	return nil, nil
// }
