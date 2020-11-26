package storage
// Copyright (C) 2020 ConsenSys Software Inc

import (
//	"log"
	"sync"
)


// Storage for key-value pairs.


// Map of key value pairs.
var keyValuePairsInstance = newKeyValueStorage()

// KeyValueStorage is the in memory implementation of storage.
type KeyValueStorage struct {
	Base
	keyValueMap     map[string]string
	keyValueMapLock sync.RWMutex
}

// Create a new instance
func newKeyValueStorage() *KeyValueStorage {
	aMap := make(map[string]string)
	aLock := sync.RWMutex{}
	var kv = KeyValueStorage{}
	kv.keyValueMap = aMap
	kv.keyValueMapLock = aLock

	var _ Storage = &kv  // Enforce interface compliance
	return &kv
}

// GetKeyValueStorage returns an instance of in-memory storage
func GetKeyValueStorage() *KeyValueStorage {
	return keyValuePairsInstance
}


// Type returns the type of storage
func (s *KeyValueStorage) Type() SType {
	return KeyValue;
}


// Put adds to the map / replace an existing value.
func (s *KeyValueStorage) Put(key, value string) {
	s.keyValueMapLock.Lock()
	s.keyValueMap[key] = value
	s.keyValueMapLock.Unlock()
}



// GetValue returns a value given a key
func (s *KeyValueStorage) GetValue(key string) (val string, exists bool) {
	val, exists = s.keyValueMap[key]
	return
}


// GetKeys returns all of the keys
func (s *KeyValueStorage) GetKeys() (keys []string) {
	keys = make([]string, len(s.keyValueMap))

	i := 0
	for k := range s.keyValueMap {
		keys[i] = k
		i++
	}
	return keys
}
