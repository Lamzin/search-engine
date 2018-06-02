package builder

type twoArrays struct {
	Key   []uint32
	Value []uint32
}

func (t twoArrays) Len() int {
	return len(t.Key)
}

func (t twoArrays) Swap(i, j int) {
	t.Key[i], t.Key[j] = t.Key[j], t.Key[i]
	t.Value[i], t.Value[j] = t.Value[j], t.Value[i]
}

func (t twoArrays) Less(i, j int) bool {
	return t.Key[i] < t.Key[j]
}

func (a *twoArrays) Merge(b *twoArrays) *twoArrays {
	c := twoArrays{
		Key:   make([]uint32, a.Len()+b.Len()),
		Value: make([]uint32, a.Len()+b.Len()),
	}
	aiter, biter := 0, 0
	for aiter+biter < c.Len() {
		if aiter == a.Len() || (biter < b.Len() && a.Key[aiter] > b.Key[biter]) {
			c.Key[aiter+biter] = b.Key[biter]
			c.Value[aiter+biter] = b.Value[biter]
			biter++
		} else {
			c.Key[aiter+biter] = a.Key[aiter]
			c.Value[aiter+biter] = a.Value[aiter]
			aiter++
		}
	}
	return &c
}
