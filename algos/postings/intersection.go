package postings

import (
	"sort"

	"github.com/lamzin/search-engine/v2/index/builder"
)

type PostingFrequencies struct {
	Posting     uint32
	Frequencies []uint32
}

func Intersect(metaArrays []*builder.MetaArrays) []uint32 {
	postToFreq := make(map[uint32][]uint32)
	for _, metaArray := range metaArrays {
		for i := 0; i < metaArray.Len(); i++ {
			posting := metaArray.Postings[i]
			postToFreq[posting] = append(postToFreq[posting], metaArray.Frequencies[i])
		}
	}

	var postingFrequenciesList []*PostingFrequencies
	for posting, frequencies := range postToFreq {
		if len(frequencies) == len(metaArrays) {
			sort.Sort(uint32Order(frequencies))
			postingFrequenciesList = append(postingFrequenciesList, &PostingFrequencies{
				Posting:     posting,
				Frequencies: frequencies,
			})
		}
	}

	sort.Sort(minMaxSorter(postingFrequenciesList))

	var postings []uint32
	for _, postingFrequencies := range postingFrequenciesList {
		postings = append(postings, postingFrequencies.Posting)
	}
	return postings
}

type minMaxSorter []*PostingFrequencies

func (m minMaxSorter) Len() int {
	return len(m)
}

func (m minMaxSorter) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m minMaxSorter) Less(i, j int) bool {
	for g := 0; g < len(m[i].Frequencies); g++ {
		if m[i].Frequencies[g] == m[j].Frequencies[g] {
			continue
		}
		return m[i].Frequencies[g] > m[j].Frequencies[g]
	}
	return m[i].Posting < m[j].Posting
}

type uint32Order []uint32

func (arr uint32Order) Len() int {
	return len(arr)
}

func (arr uint32Order) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (arr uint32Order) Less(i, j int) bool {
	return arr[i] < arr[j]
}
