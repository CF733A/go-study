package main

import (
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"time"

	"github.com/sigurn/crc16"
	"github.com/sigurn/crc8"
)

type HashMap struct {
	size     int
	buckets  []*bucket
	hashFunc func(string) uint64
}

type bucket struct {
	key   string
	value interface{}
	next  *bucket
}

type HashOption func(*HashMap)

func NewHashMap(size int, opts ...HashOption) *HashMap {
	hm := &HashMap{
		size:    size,
		buckets: make([]*bucket, size),
		hashFunc: func(s string) uint64 {
			return uint64(crc32.ChecksumIEEE([]byte(s)))
		},
	}
	for _, opt := range opts {
		opt(hm)
	}
	return hm
}

func WithHashCRC8() HashOption {
	table := crc8.MakeTable(crc8.CRC8)
	return func(hm *HashMap) {
		hm.hashFunc = func(s string) uint64 {
			return uint64(crc8.Checksum([]byte(s), table))
		}
	}
}

func WithHashCRC16() HashOption {
	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
	return func(hm *HashMap) {
		hm.hashFunc = func(s string) uint64 {
			return uint64(crc16.Checksum([]byte(s), table))
		}
	}
}

func WithHashCRC32() HashOption {
	return func(hm *HashMap) {
		hm.hashFunc = func(s string) uint64 {
			return uint64(crc32.ChecksumIEEE([]byte(s)))
		}
	}
}

func WithHashCRC64() HashOption {
	table := crc64.MakeTable(crc64.ECMA)
	return func(hm *HashMap) {
		hm.hashFunc = func(s string) uint64 {
			return uint64(crc64.Checksum([]byte(s), table))
		}
	}
}

type HashMaper interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
}

func (h *HashMap) Set(key string, value interface{}) {
	hash := h.hashFunc(key)
	index := hash % uint64(h.size)

	if h.buckets[index] == nil {
		h.buckets[index] = &bucket{key: key, value: value}
		return
	}

	current := h.buckets[index]
	for {
		if current.key == key {
			current.value = value
			return
		}

		if current.next == nil {
			current.next = &bucket{key: key, value: value}
			return
		}
		current = current.next
	}
}

func (h *HashMap) Get(key string) (interface{}, bool) {
	hash := h.hashFunc(key)
	index := hash % uint64(h.size)

	current := h.buckets[index]
	for current != nil {
		if current.key == key {
			return current.value, true
		}
		current = current.next
	}
	return nil, false
}

func MeassureTime(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}

func main() {
	m := NewHashMap(16, WithHashCRC64())
	since := MeassureTime(func() {
		m.Set("key", "value")

		if value, ok := m.Get("key"); ok {
			fmt.Println(value)
		}
	})
	fmt.Println(since)

	m = NewHashMap(16, WithHashCRC32())
	since = MeassureTime(func() {
		m.Set("key", "value")

		if value, ok := m.Get("key"); ok {
			fmt.Println(value)
		}
	})
	fmt.Println(since)

	m = NewHashMap(16, WithHashCRC16())
	since = MeassureTime(func() {
		m.Set("key", "value")

		if value, ok := m.Get("key"); ok {
			fmt.Println(value)
		}
	})
	fmt.Println(since)

	m = NewHashMap(16, WithHashCRC8())
	since = MeassureTime(func() {
		m.Set("key", "value")

		if value, ok := m.Get("key"); ok {
			fmt.Println(value)
		}
	})
	fmt.Println(since)
}
