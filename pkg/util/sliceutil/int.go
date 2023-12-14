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


//RemoveIf 移除符合数组中某个条件的元素
func (m IntSlice)RemoveIf(condition func(int) bool) []int {
	if m == nil {
		return nil
	}
	// 第一次迭代：标记要删除的元素
	marked := make([]bool, len(m))
	for i, num := range m {
		if condition(num) {
			marked[i] = true
		}
	}

	// 第二次迭代：删除标记为 true 的元素
	result := make([]int, 0, len(m))
	for i, num := range m {
		if !marked[i] {
			result = append(result, num)
		}
	}

	return result
}

//AppendIf 追加符合莫格条件的元素到数组
func (m IntSlice)AppendIf(condition func(int) bool,sl []int) []int {
	if sl == nil {
		return m
	}

	result := make([]int, 0, len(m))
	for _, str := range m {
		result =append(result,str)
	}

	for _,s:=range sl{
		if condition(s) {
			result = append(result,s)
		}
	}
	return result
}
