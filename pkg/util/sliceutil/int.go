package sliceutil

type IntSlice []int

func (m IntSlice) DeepCopy() IntSlice {
	o := make([]int, 0, len(m))
	o = append(o, m...)
	return o
}

func (m IntSlice) Append(target ...int) IntSlice {
	if m == nil {
		o := make([]int, 0, len(target))
		return append(o, target...)
	}

	return append(m, target...)
}
