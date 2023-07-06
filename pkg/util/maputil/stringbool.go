package maputil

type StringBoolMap map[string]bool

func (m StringBoolMap) DeepCopy() StringBoolMap {
	o := make(map[string]bool, len(m))
	for k, v := range m {
		o[k] = v
	}
	return o
}

func (m StringBoolMap) Init() StringBoolMap {
	if m == nil {
		return make(map[string]bool)
	}
	return m
}

func (m StringBoolMap) Delete(key string) {
	if m == nil {
		return
	}
	delete(m, key)
}

func (m StringBoolMap) Has(key string) bool {
	if m != nil {
		if _, exist := m[key]; exist {
			return true
		}
	}
	return false
}

func (m StringBoolMap) Set(key string, value bool) StringBoolMap {
	if m == nil {
		o := make(map[string]bool)
		o[key] = value
		return o
	}
	m[key] = value
	return m
}

func (m StringBoolMap) Get(key string) bool {
	if m == nil {
		return false
	}
	v, _ := m[key]
	return v
}
