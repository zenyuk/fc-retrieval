package storage

// Copyright (C) 2020 ConsenSys Software Inc

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Storage for key-value pairs using mongodb.

// Name of mongodb host.
const mongoHostname = "mongodb"
const mongoURI = "mongodb://" + mongoHostname

const dbName = "kv"
const dbCollection = "values"

// Lazy initialize instance variable.
var mongoInstance *MongoStorage

// MongoStorage is the Mongo db implementation of storage.
type MongoStorage struct {
	Base
	session *mgo.Session
}



// Create a new instance
func newMongoStorage() *MongoStorage {
	m := MongoStorage{}
	m.session = getSession()

	var _ Storage = &m  // Enforce interface compliance

	return &m
}

// GetMongoStorage returns the Mongo db implementation of the KeyValueStorage.
func GetMongoStorage() *MongoStorage {
	if mongoInstance == nil {
		mongoInstance =  newMongoStorage()
	}
	return mongoInstance
}

// Type returns the type of storage
func (s *MongoStorage) Type() SType {
	return Mongo;
}



// Put adds to the map / replace an existing value.
func (s *MongoStorage) Put(key, value string) {
	session := s.session.Copy()
	defer session.Close()

	kv := CollectionKeyValue{}
	kv.Key = key
	kv.Value = value

	collection := session.DB(dbName).C("values")

	err := collection.Insert(kv)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
}


// GetValue returns a value given a key
func (s *MongoStorage) GetValue(key string) (val string, exists bool) {
	session := s.session.Copy()
	defer session.Close()

	kv := CollectionKeyValue{}

	collection := session.DB(dbName).C(dbCollection)
	err := collection.Find(bson.M{"Key": key}).One(&kv)
	if err != nil {
		return "", false
	}

	return kv.Value, true
}


// GetKeys returns all of the keys
func (s *MongoStorage) GetKeys() (keys []string) {
	session := s.session.Copy()
	defer session.Close()

	kv := CollectionKeyValue{}

	collection := session.DB(dbName).C(dbCollection)
	iter := collection.Find(bson.M{}).Iter()

	keys = make([]string, 0)   // Even if there are no elements, return something
	for iter.Next(&kv) {
		keys = append(keys, kv.Key)
	}
	return
}



// CollectionKeyValue holds the database collection.
type CollectionKeyValue struct {
	Key   	string	`bson:"Key"`
	Value 	string  `bson:"Value"`
}


func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial(mongoURI)

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	// TODO what does this do?
	// s.SetSafe(&mgo.Safe{})
	return s
}

