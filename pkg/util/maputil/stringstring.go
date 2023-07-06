package maputil

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
