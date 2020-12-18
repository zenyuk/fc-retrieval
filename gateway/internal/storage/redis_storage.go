package storage
// Copyright (C) 2020 ConsenSys Software Inc

import (
	"github.com/xyproto/simpleredis"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)


// Storage for key-value pairs using mongodb.


// Name of mongodb host.
const redisHostname = "redis:6379"
//const REDIS_HOSTNAME = "redis-master:6379"

const kvName = "KV"

// Lazy initialize instance variable.
var redisInstance *RedisStorage

// RedisStorage holds the Redis implementation of the storage
type RedisStorage struct {
	Base
	masterPool *simpleredis.ConnectionPool
	//slavePool  *simpleredis.ConnectionPool
	kv *simpleredis.KeyValue
}



// Create a new instance
func newRedisStorage() *RedisStorage {
	m := RedisStorage{}
	m.masterPool = simpleredis.NewConnectionPoolHost(redisHostname)
	m.kv = simpleredis.NewKeyValue(m.masterPool, kvName)
	m.kv.SelectDatabase(1)
	var _ Storage = &m  // Enforce interface compliance
	return &m
}

// GetRedisStorage returns the Redis implementation of KeyValueStorage.
func GetRedisStorage() *RedisStorage {
	if redisInstance == nil {
		redisInstance =  newRedisStorage()
	}
	return redisInstance
}


// Type returns the type of storage.
func (s *RedisStorage) Type() SType {
	return Redis;
}



// Put adds to the map / replace an existing value.
func (s *RedisStorage) Put(key, value string) {
	err := s.kv.Set(key, value)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
}


// GetValue returns a value given a key
func (s *RedisStorage) GetValue(key string) (val string, exists bool) {
	val, err := s.kv.Get(key)
	if err != nil {
		return "", false
	}

	return val, true
}


// GetKeys returns all of the keys
func (s *RedisStorage) GetKeys() (keys []string) {
	keys = make([]string, 0)   // Even if there are no elements, return something
	keys = append(keys, "Not implemented yet")
	return
}


