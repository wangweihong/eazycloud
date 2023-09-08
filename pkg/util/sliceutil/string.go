package sliceutil

import (
	"sort"
	"strings"
)

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

//GetRepeat find slice repeat data and repeat num
func (m StringSlice) GetRepeat() (map[string]int, bool) {
	if m != nil {
		var r map[string]int
		s := make(map[string]struct{})
		for _, v := range m {
			if _, exist := s[v]; exist {
				if r == nil {
					r = make(map[string]int)
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
func (m StringSlice) SortDesc() []string {
	if m != nil {
		sort.Slice(m, func(i, j int) bool {
			return m[i] > m[j]
		})
		return m
	}

	return nil
}

// Sort Ascascending sort
func (m StringSlice) SortAsc() []string {
	if m != nil {
		sort.Slice(m, func(i, j int) bool {
			return m[i] < m[j]
		})
		return m
	}

	return nil
}

func (m StringSlice) HasEmpty() (int, bool) {
	if m != nil {
		var eN int
		for _, v := range m {
			if v == "" {
				eN++
			}
		}
		return eN, eN != 0
	}

	return 0, false
}

func (m StringSlice) Cut(data string) []string {
	var index int = -1
	for i, v := range m {
		if v == data {
			index = i
			break
		}
	}

	if index == -1 {
		return m
	}

	return append(m[:index], m[index+1:]...)
}

func (m StringSlice) FallBehind(data string) []string {
	n := m.Cut(data)
	return append(n, data)
}

func (m StringSlice) TrimSpace() []string {
	if m == nil {
		return nil
	}

	n := make([]string, 0, len(m))
	for _, v := range m {
		if strings.TrimSpace(v) != "" {
			n = append(n, v)
		}
	}
	return n
}
