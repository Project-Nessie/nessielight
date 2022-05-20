package tgolf

import (
	"strings"
	"sync"
)

// simple abstraction for key-value database
type KVDatabase interface {
	// get value based on key
	Get(key string) interface{}
	Set(key string, value interface{}) error
	// return a subdatabase with prefix. e. g. subdb.get(key) is equal to
	// db.get(prefix + key), ant note that subdb & db share the same storage
	Sub(prefix string) KVDatabase
	// traverse all data
	ForEach(handler func(key string, value interface{}))
}

// MemoryDB implememts KVDatabase, which is a simple storage based on runtime memory
type MemoryDB struct {
	// 多线程读写锁
	*sync.RWMutex
	db map[string]interface{}
	// 当前存储的前缀
	prefix string
}

var _ KVDatabase = (*MemoryDB)(nil)

func (r *MemoryDB) Get(key string) interface{} {
	r.RLock()
	defer r.RUnlock()
	return r.db[r.prefix+key]
}

func (r *MemoryDB) Set(key string, value interface{}) error {
	r.Lock()
	defer r.Unlock()
	if value == nil {
		delete(r.db, r.prefix+key)
	} else {
		r.db[r.prefix+key] = value
	}
	return nil
}

func (r *MemoryDB) ForEach(handler func(key string, value interface{})) {
	for k, v := range r.db {
		if strings.HasPrefix(k, r.prefix) {
			key := k[len(r.prefix):]
			handler(key, v)
		}
	}
}

func (r *MemoryDB) Sub(prefix string) KVDatabase {
	db := *r
	db.prefix += prefix
	return &db
}

func NewMemoryDB() MemoryDB {
	return MemoryDB{db: make(map[string]interface{}), RWMutex: &sync.RWMutex{}}
}
