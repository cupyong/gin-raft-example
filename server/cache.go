package server

import (
	"sync"
	"encoding/json"
	"io"
)

type Cache struct {
	data map[string]string
	sync.RWMutex
}
func NewCache() *Cache {
	cache := &Cache{}
	cache.data = make(map[string]string)
	return cache
}

func (c *Cache) Get(key string) string {
	c.RLock()
	ret := c.data[key]
	c.RUnlock()
	return ret
}

func (c *Cache) Set(key string, value string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = value
	return nil
}

// Marshal serializes cache data
func (c *Cache) Marshal() ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	dataBytes, err := json.Marshal(c.data)
	return dataBytes, err
}

// UnMarshal deserializes cache data
func (c *Cache) UnMarshal(serialized io.ReadCloser) error {
	var newData map[string]string
	if err := json.NewDecoder(serialized).Decode(&newData); err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()
	c.data = newData

	return nil
}
