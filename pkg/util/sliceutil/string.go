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
