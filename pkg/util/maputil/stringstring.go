package maputil

import "github.com/wangweihong/eazycloud/pkg/sets"

type StringStringMap map[string]string

func (m StringStringMap) DeepCopy() StringStringMap {
	o := make(map[string]string, len(m))
	for k, v := range m {
		o[k] = v
	}
	return o
}

func (m StringStringMap) Init() StringStringMap {
	if m == nil {
		return make(map[string]string)
	}
	return m
}

func (m StringStringMap) Delete(key string) {
	if m == nil {
		return
	}
	delete(m, key)
}


func (m StringStringMap) DeleteIfKey(condition func(string) bool) {
	if m == nil {
		return
	}

	for k := range m {
		if condition(k) {
			delete(m, k)
		}
	}
}

func (m StringStringMap) DeleteIfValue(condition func(string) bool) {
	if m == nil {
		return
	}

	for k,v := range m {
		if condition(v) {
			delete(m, k)
		}
	}
}

func (m StringStringMap) Has(key string) bool {
	if m != nil {
		if _, exist := m[key]; exist {
			return true
		}
	}
	return false
}

func (m StringStringMap) Set(key string, value string) StringStringMap {
	if m == nil {
		o := make(map[string]string)
		o[key] = value
		return o
	}
	m[key] = value
	return m
}

func (m StringStringMap) Get(key string) string {
	if m == nil {
		return ""
	}
	v, _ := m[key]
	return v
}

func (m StringStringMap) Keys() []string {
	if m == nil {
		return []string{}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m StringStringMap) ToSetString() sets.String {
	ss := sets.NewString()

	if m == nil {
		return ss
	}
	for k:=range m{
		ss.Insert(k)
	}
	return ss
}