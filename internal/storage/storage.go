package storage
// Copyright (C) 2020 ConsenSys Software Inc


// Storage is the interface which all storage providers implement.
type Storage interface {
	Type() SType
	Put(key, value string)
	GetValue(key string) (val string, exists bool)
	GetKeys() (keys []string)
}


// Base is the base class for all storage providers.
type Base struct {
	// Currently, nothing in base class.
}




// SType is an enum to define type of storage.
type SType int

// Values of types of storage implementation
const (
	KeyValue SType = iota //1
	Mongo // 2
	Redis // 3
)

// String returns the type of storage as a string
func (t SType) String() string {
	switch t {
	case KeyValue:
		return "Key Value"
	case Mongo:
		return "Mongo"
	case Redis:
		return "Redis"

	default:
		panic("unrecognized storage type")
	}
}


// GetSingleInstance is a factory method to get the single instance of the storage.
func GetSingleInstance(t SType) *Storage {
	var s Storage
	switch t {
	case KeyValue:
		s = GetKeyValueStorage()
	case Mongo:
		s = GetMongoStorage()
	case Redis:
		s = GetRedisStorage()
	default:
		panic("unrecognized storage type")
	}
	return &s
}



