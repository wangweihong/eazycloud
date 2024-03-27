package sequential

import "sync"

// 顺序表, 结合表的功能，提供插入数据的顺序
type Map struct {
	lock sync.RWMutex
	data []interface{}
	//用于记录插入的键的顺序
	key     []interface{}
	indices map[interface{}]int
}

func NewSequentialMap() *Map {
	return &Map{
		data:    make([]interface{}, 0),
		key:     make([]interface{}, 0),
		indices: make(map[interface{}]int),
	}
}

func (m *Map) Get(value interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if value != nil {
		i, exist := m.indices[value]
		if exist {
			return m.data[i]
		}
	}
	return nil
}

func (m *Map) Has(key interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if key == nil {
		return false
	}

	_, exist := m.indices[key]
	return exist
}

func (m *Map) ForEach(f func(value interface{}) error) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, v := range m.data {
		err := f(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Map) Inject(key interface{}, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if key == nil || value == nil {
		return
	}

	i, exist := m.indices[key]
	if exist {
		// update new value
		m.data[i] = value

		return
	}

	m.data = append(m.data, value)
	m.key = append(m.key, key)
	m.indices[key] = len(m.data) - 1
}

func (m *Map) Map() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	nm := make(map[interface{}]interface{})
	for k, v := range m.indices {
		nm[k] = m.data[v]
	}
	return nm
}

func (m *Map) List() []interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	nl := make([]interface{}, 0, len(m.data))
	for _, v := range m.data {
		nl = append(nl, v)
	}
	return nl
}

func (m *Map) Keys() []interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	keys := make([]interface{}, 0, len(m.data))
	keys = append(keys, m.key...)
	return keys
}

func (m *Map) Values() []interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	vals := make([]interface{}, 0, len(m.data))
	vals = append(vals, m.data...)
	return vals
}

func (m *Map) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return len(m.data)
}

func (m *Map) Delete(key interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if key == nil {
		return
	}

	index, ok := m.indices[key]
	if !ok {
		return
	}

	delete(m.indices, key)

	m.data = append(m.data[:index], m.data[index+1:]...)
	m.key = append(m.key[:index], m.key[index+1:]...)
	// 更新被移除元素后面元素的索引
	for i := index; i < len(m.data); i++ {
		m.indices[m.data[i]] = i
	}
}

func (m *Map) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.data = make([]interface{}, 0)
	m.key = make([]interface{}, 0)
	m.indices = make(map[interface{}]int)
}
