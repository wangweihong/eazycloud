package sliceutil

import "sort"

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

//GetRepeat find slice repeat data and repeat num
func (m IntSlice) GetRepeat() (map[int]int, bool) {
	if m != nil {
		var r map[int]int
		s := make(map[int]struct{})
		for _, v := range m {
			if _, exist := s[v]; exist {
				if r == nil {
					r = make(map[int]int)
				}
				num, _ := r[v]
				if num == 0 {
					num = 1
				}
				num++
				r[v] = num
			}
			s[v] = struct{}{}
		}

		return r, !(len(r) == 0)
	}

	return nil, false
}

//SortDesc Descending sort
func (m IntSlice) SortDesc() []int {
	if m != nil {
		sort.Slice(m, func(i, j int) bool {
			return m[i] > m[j]
		})
		return m
	}

	return nil
}

// Sort Ascascending sort
func (m IntSlice) SortAsc() []int {
	if m != nil {
		sort.Slice(m, func(i, j int) bool {
			return m[i] < m[j]
		})
		return m
	}

	return nil
}
