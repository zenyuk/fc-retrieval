package api
// Copyright (C) 2020 ConsenSys Software Inc

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/storage"
)

// KeyValue holds keys and values that are stored
type KeyValue struct {
	Key   string
	Value string
}




func setValue(w rest.ResponseWriter, r *rest.Request) {
	log.Println("setValue: start")
//TODO this doesn't compile with new version of Go	log.Printf("setValue: request: \n%s\n", r)
	key, value, e := extractKeyValue(w, r)
	if e {
		log.Println("setValue: Extract key value pair error")
		// NOTE: HTTP error already set-up.
		return
	}
	if !isB64OrSimpleAscii(key) {
		log.Println("setValue: key not Base64Url or ASCII encoded")
		rest.Error(w, "key not Base64Url or ASCII encoded", http.StatusBadRequest)
		return
	}
	if !isB64OrSimpleAscii(value) {
		log.Println("setValue: value not Base64Url or ASCII encoded")
		rest.Error(w, "value not Base64Url or ASCII encoded", http.StatusBadRequest)
		return
	}

	store := getStorage();
	//store := storage.GetKeyValueStorage()

	store.Put(key, value)

	w.WriteHeader(http.StatusOK)

	log.Println("setValue: done")
}




func getKeyValues(w rest.ResponseWriter, r *rest.Request) {
	store := getStorage();

	queryValues := r.URL.Query()
	numQueryKeyValuePairs := len(queryValues)

	key := queryValues.Get("Key")
	log.Printf("Key: %s\n", key)
	if len(key) != 0 {
		if numQueryKeyValuePairs != 1 {
			rest.Error(w, "Invalid parameter. Parameter values: Key=key, or no parameters", http.StatusBadRequest)
			return
		}

		if !isB64OrSimpleAscii(key) {
			rest.Error(w, "key not Base64Url or ASCII encoded", http.StatusBadRequest)
			return
		}

		value, exists := store.GetValue(key)
		//log.Printf("value: %s\n", value)
		if (exists) {
			w.WriteJson(&value)
			return
		}
		value = "ERROR: Key value not set"
		w.WriteJson(&value)
		return
	}
	// return all keys.

	if numQueryKeyValuePairs != 0 {
		rest.Error(w, "Invalid parameter. Parameter values: Key=key, or no parameters", http.StatusBadRequest)
		return
	}

	keys := store.GetKeys()
	w.WriteJson(&keys)
}


func getStorage() storage.Storage {
	store := storage.GetSingleInstance(storage.Redis)
	return *store
}


func extractKeyValue(w rest.ResponseWriter, r *rest.Request) (k, v string, e bool){
	log.Println("ExtractKeyValue: start")

	e = true

	keyValue := KeyValue{}

	err := r.DecodeJsonPayload(&keyValue)
	if err != nil {
		log.Println("Error: ", err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Key: %s, Value: %s\n", keyValue.Key, keyValue.Value)
	k = keyValue.Key
	v = keyValue.Value

	log.Println("ExtractKeyValue: done")

	e = false
	return
}
