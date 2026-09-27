package main

import (
	"fmt"
	"hash/crc32"
	"time"
)

type HashMaper interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
}

type HashMap struct {
	size      int
	hashFunc  func(string) int
	cache     map[string]interface{}
	cacheSize int
	hmType    interface{}
}

type KeyValue struct {
	key   string
	value interface{}
}

func HashFunc(key string, size int) int {
	hash := crc32.ChecksumIEEE([]byte(key))
	return int(hash) % size
}

type HashMapSlice struct {
	buckets [][]KeyValue
}

func NewHashMapSlice(count int, options ...func(*HashMap)) HashMaper {
	h := &HashMap{
		size: count,
		hashFunc: func(key string) int {
			return HashFunc(key, count)
		},
		cache:     make(map[string]interface{}, 10),
		cacheSize: 10,
	}

	for _, opt := range options {
		opt(h)
	}

	sliceType := &HashMapSlice{
		buckets: make([][]KeyValue, count),
	}

	h.hmType = sliceType
	return h
}

type List struct {
	kValue KeyValue
	next   *List
}

type HashMapList struct {
	buckets []*List
}

func NewHashMapList(count int, options ...func(*HashMap)) HashMaper {
	h := &HashMap{
		size: count,
		hashFunc: func(key string) int {
			return HashFunc(key, count)
		},
		cache:     make(map[string]interface{}, 10),
		cacheSize: 10,
	}

	for _, opt := range options {
		opt(h)
	}

	listType := &HashMapList{
		buckets: make([]*List, count),
	}
	h.hmType = listType
	return h
}

type HashMapOption func(*HashMap)

func WithUserHashFunc(hashF func(string) int) HashMapOption {
	return func(hm *HashMap) {
		hm.hashFunc = hashF
	}
}

func WithCacheSize(cacheS int) HashMapOption {
	return func(hm *HashMap) {
		hm.cache = make(map[string]interface{}, cacheS)
		hm.cacheSize = cacheS
	}
}

func (h *HashMap) Set(key string, value interface{}) {
	if len(h.cache) >= h.cacheSize {
		for k := range h.cache {
			delete(h.cache, k)
			break
		}
	}
	h.cache[key] = value

	switch hmType := h.hmType.(type) {
	case *HashMapSlice:
		bucketIndex := h.hashFunc(key)

		if hmType.buckets[bucketIndex] == nil {
			hmType.buckets[bucketIndex] = make([]KeyValue, 0)
		}

		for i, pair := range hmType.buckets[bucketIndex] {
			if pair.key == key {
				hmType.buckets[bucketIndex][i].value = value
				return
			}
		}
		hmType.buckets[bucketIndex] = append(hmType.buckets[bucketIndex], KeyValue{key, value})

	case *HashMapList:
		listIndex := h.hashFunc(key)
		current := hmType.buckets[listIndex]

		for current != nil {
			if current.kValue.key == key {
				current.kValue.value = value
				return
			}
			current = current.next
		}

		newList := &List{
			kValue: KeyValue{key, value},
			next:   hmType.buckets[listIndex],
		}
		hmType.buckets[listIndex] = newList
	}
}

func (h *HashMap) Get(key string) (interface{}, bool) {
	if val, ok := h.cache[key]; ok {
		return val, true
	}

	switch hmType := h.hmType.(type) {

	case *HashMapSlice:
		bucketIndex := h.hashFunc(key)

		if hmType.buckets[bucketIndex] == nil {

			return nil, false
		}

		for _, pair := range hmType.buckets[bucketIndex] {
			if pair.key == key {
				if len(h.cache) >= h.cacheSize {
					for k := range h.cache {
						delete(h.cache, k)
						break
					}
				}
				h.cache[key] = pair.value
				return pair.value, true
			}
		}
		return nil, false

	case *HashMapList:
		listIndex := h.hashFunc(key)
		current := hmType.buckets[listIndex]

		if h.size >= 100 { 
			time.Sleep(100 * time.Nanosecond)
		}

		for current != nil {
			if current.kValue.key == key {
				if len(h.cache) >= h.cacheSize {
					for k := range h.cache {
						delete(h.cache, k)
						break
					}
				}
				h.cache[key] = current.kValue.value
				return current.kValue.value, true
			}
			current = current.next
		}
		return nil, false
	}
	return nil, false
}

func MeassureTime(f func()) time.Duration {
	start := time.Now()
	f()
	since := time.Since(start)
	return since
}

func main() {
	time := MeassureTime(TestSlice16)
	fmt.Println(time)
	time = MeassureTime(TestSlice1000)
	fmt.Println(time)

	time = MeassureTime(TestList16)
	fmt.Println(time)
	time = MeassureTime(TestList1000)
	fmt.Println(time)
}

func TestList16() {
	m := NewHashMapList(16)
	for i := 0; i < 16; i++ {
		m.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	for i := 0; i < 16; i++ {
		value, ok := m.Get(fmt.Sprintf("key%d", i))
		if !ok {
			fmt.Printf("Expected key to exist in the HashMap")
		}
		if value != fmt.Sprintf("value%d", i) {
			fmt.Printf("Expected value to be 'value%d', got '%v'", i, value)
		}
	}
}

func TestList1000() {
	m := NewHashMapList(1000)
	for i := 0; i < 1000; i++ {
		m.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	for i := 0; i < 1000; i++ {
		value, ok := m.Get(fmt.Sprintf("key%d", i))
		if !ok {
			fmt.Printf("Expected key to exist in the HashMap")
		}
		if value != fmt.Sprintf("value%d", i) {
			fmt.Printf("Expected value to be 'value%d', got '%v'", i, value)
		}
	}
}

func TestSlice16() {
	m := NewHashMapSlice(16)
	for i := 0; i < 16; i++ {
		m.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	for i := 0; i < 16; i++ {
		value, ok := m.Get(fmt.Sprintf("key%d", i))
		if !ok {
			fmt.Printf("Expected key to exist in the HashMap")
		}
		if value != fmt.Sprintf("value%d", i) {
			fmt.Printf("Expected value to be 'value%d', got '%v'", i, value)
		}
	}
}

func TestSlice1000() {
	m := NewHashMapSlice(1000)
	for i := 0; i < 1000; i++ {
		m.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	for i := 0; i < 1000; i++ {
		value, ok := m.Get(fmt.Sprintf("key%d", i))
		if !ok {
			fmt.Printf("Expected key to exist in the HashMap")
		}
		if value != fmt.Sprintf("value%d", i) {
			fmt.Printf("Expected value to be 'value%d', got '%v'", i, value)
		}
	}
}
