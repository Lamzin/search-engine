package builder

type Meta struct {
	Posting   uint32
	Frequency uint32
}

type MetaArrays struct {
	Postings    []uint32
	Frequencies []uint32
}

func (a MetaArrays) Len() int {
	return len(a.Postings)
}

func (a MetaArrays) Swap(i, j int) {
	a.Postings[i], a.Postings[j] = a.Postings[j], a.Postings[i]
	a.Frequencies[i], a.Frequencies[j] = a.Frequencies[j], a.Frequencies[i]
}

func (a MetaArrays) Less(i, j int) bool {
	return a.Postings[i] < a.Postings[j]
}

func (a *MetaArrays) Merge(b *MetaArrays) *MetaArrays {
	c := MetaArrays{
		Postings:    make([]uint32, a.Len()+b.Len()),
		Frequencies: make([]uint32, a.Len()+b.Len()),
	}
	aiter, biter := 0, 0
	for aiter+biter < c.Len() {
		if aiter == a.Len() || (biter < b.Len() && a.Postings[aiter] > b.Postings[biter]) {
			c.Postings[aiter+biter] = b.Postings[biter]
			c.Frequencies[aiter+biter] = b.Frequencies[biter]
			biter++
		} else {
			c.Postings[aiter+biter] = a.Postings[aiter]
			c.Frequencies[aiter+biter] = a.Frequencies[aiter]
			aiter++
		}
	}
	return &c
}
