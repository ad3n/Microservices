package main

import (
	"fmt"
	"time"

	"github.com/gadelkareem/cachita"
)

type cache struct {
	pool cachita.Cache
}

func NewCache() *cache {
	return &cache{
		pool: cachita.Memory(),
	}
}

func (c cache) set(key string, value *data, lifetime int) {
	err := c.pool.Put(key, value, time.Duration(lifetime)*time.Second)
	if err != nil {
		fmt.Println(err)
	}
}

func (c cache) get(key string) (*data, bool) {
	data := data{}
	err := c.pool.Get(key, &data)
	if err != nil {
		return nil, false
	}

	return &data, true
}

func (c cache) invalidate(key string) {
	c.pool.Invalidate(key)
}
