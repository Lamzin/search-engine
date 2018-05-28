package postings

import (
	"sort"
)

type PostingsFrequency struct {
	DocID     int
	Frequency int
}

func Intersect(postingLists [][]int) []PostingsFrequency {
	docToFrequency := make(map[int]int)
	for _, postingList := range postingLists {
		for _, posting := range postingList {
			docToFrequency[posting]++
		}
	}

	frequencyList := make([]PostingsFrequency, len(docToFrequency))
	i := 0
	for docID, frequency := range docToFrequency {
		frequencyList[i].DocID = docID
		frequencyList[i].Frequency = frequency
		i++
	}

	sort.Sort(ByFrequency(frequencyList))
	return frequencyList
}

type ByFrequency []PostingsFrequency

func (arr ByFrequency) Len() int {
	return len(arr)
}

func (arr ByFrequency) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (arr ByFrequency) Less(i, j int) bool {
	if arr[i].Frequency == arr[j].Frequency {
		return arr[i].DocID < arr[j].DocID
	}
	return arr[i].Frequency > arr[j].Frequency
}
