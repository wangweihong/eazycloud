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

//HasRepeat slice has repeated data
func (m IntSlice) HasRepeat() bool {
	if m != nil {
		s := make(map[int]struct{})
		for _, v := range m {
			if _, exist := s[v]; exist {
				return true
			}
			s[v] = struct{}{}
		}
	}

	return false
}
