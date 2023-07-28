package sliceutil

type StringSlice []string

func (m StringSlice) DeepCopy() StringSlice {
	o := make([]string, 0, len(m))
	o = append(o, m...)
	return o
}

func (m StringSlice) Append(target ...string) StringSlice {
	if m == nil {
		o := make([]string, 0, len(target))
		return append(o, target...)
	}

	return append(m, target...)
}

//HasRepeat slice has repeated data
func (m StringSlice) HasRepeat() bool {
	if m != nil {
		s := make(map[string]struct{})
		for _, v := range m {
			if _, exist := s[v]; exist {
				return true
			}
			s[v] = struct{}{}
		}
	}

	return false
}
