package datasource

import (
	"errors"
	. "httpstore/logging"
	"sync"
)

// Errors
var (
	ErrKeyNotFound          = errors.New("key not found")
	ErrKvStoreAlreadyExists = errors.New("key value store has already been initialized")
	ErrKvStoreDoesNotExist  = errors.New("key value store has not been initialized")
	ErrValueUpdateForbidden = errors.New("value update not allowed, wrong owner")
	ErrValueDeleteForbidden = errors.New("value delete not allowed, wrong owner")
)

// Datasource and corresponding methods
type Datasource struct {
	kvStore map[Key]Data
	mutex   sync.RWMutex
}

func NewDatasource() Datasource {
	InfoLogger.Println("Created New Datasource with keystore")
	kvStore := make(map[Key]Data)
	mutex := sync.RWMutex{}
	return Datasource{kvStore, mutex}
}

func (ds *Datasource) isOpen() bool {
	return ds != nil
}

func (ds *Datasource) isClosed() bool {
	return ds == nil
}

func (ds *Datasource) Lock() {
	ds.mutex.Lock()
}
func (ds *Datasource) Unlock() {
	ds.mutex.Unlock()
}
func (ds *Datasource) RLock() {
	ds.mutex.RLock()
}
func (ds *Datasource) RUnlock() {
	ds.mutex.RUnlock()
}

func (ds *Datasource) Size() int {
	if ds.isClosed() {
		WarningLogger.Println("Did not get actual size because key value store is nil")
		return 0
	}
	return len(ds.kvStore)
}

func (ds *Datasource) EmptyKvStore() error {
	if ds.isClosed() {
		ErrorLogger.Println("Cannot make key value store nil when already nil")
		return ErrKvStoreDoesNotExist
	} else {
		InfoLogger.Println("Removed key store from datasource")
		ds.kvStore = nil
	}
	return nil
}

func (ds *Datasource) CreateKvStore() error {
	if ds.isOpen() {
		ErrorLogger.Println("Cannot instantiate key value store when already instantiated")
		return ErrKvStoreAlreadyExists
	} else {
		InfoLogger.Println("Created new key store for datasource")
		ds.kvStore = map[Key]Data{}
	}
	return nil
}

func (ds *Datasource) Put(key Key, newData Data) error {
	if ds.isClosed() {
		ErrorLogger.Printf("Cannot insert %s because key value store is nil", key)
		return ErrKvStoreDoesNotExist
	}
	ds.Lock()
	defer ds.Unlock()
	existingData, ok := ds.kvStore[key]
	if ok && !Authorized(&existingData, newData.owner) {
		WarningLogger.Printf("Cannot update %s because owners do not match", key)
		return ErrValueUpdateForbidden
	} else {
		ds.kvStore[key] = newData
	}
	return nil
}

func (ds *Datasource) Contains(key Key) (bool, error) {
	if ds.isClosed() {
		ErrorLogger.Printf("Cannot check if %s exists because key value store is nil", key)
		return false, ErrKvStoreDoesNotExist
	}
	ds.RLock()
	defer ds.RUnlock()
	_, ok := ds.kvStore[key]
	return ok, nil
}

func (ds *Datasource) Get(key Key) (*Data, error) {
	if ds.isClosed() {
		ErrorLogger.Printf("Cannot get %s because key value store is nil", key)
		return nil, ErrKvStoreDoesNotExist
	}
	ds.RLock()
	defer ds.RUnlock()
	existingData, ok := ds.kvStore[key]
	if !ok {
		InfoLogger.Printf("Cannot get %s because it does not exist", key)
		return nil, ErrKeyNotFound
	}
	return &existingData, nil
}

func (ds *Datasource) Delete(key Key, user string) error {
	if ds.isClosed() {
		ErrorLogger.Printf("Cannot delete %s because key value store is nil", key)
		return ErrKvStoreDoesNotExist
	}
	ds.Lock()
	defer ds.Unlock()
	existingData, ok := ds.kvStore[key]
	if !ok {
		ErrorLogger.Printf("Cannot delete %s because it does not exist", key)
		return ErrKeyNotFound
	} else if Authorized(&existingData, user) {
		delete(ds.kvStore, key)
		InfoLogger.Printf("Deleted %s with value <%s>", key, existingData.value)
		return nil
	} else {
		WarningLogger.Printf("Cannot update %s because owners do not match", key)
		return ErrValueDeleteForbidden
	}
}

// GetAllEntries Generate and return all datasource entries
func (ds *Datasource) GetAllEntries() []Entry {
	//TODO May need to optimize for larger key store sets i.e. fan in fan out retrieval
	ds.Lock()
	defer ds.Unlock()
	entries := make([]Entry, ds.Size())
	i := 0
	for key, data := range ds.kvStore {
		entries[i] = NewEntry(key, data)
		i++
	}
	return entries
}

// GetEntry Generate and return a single datasource entry
func (ds *Datasource) GetEntry(key Key) (Entry, error) {
	if ds.isClosed() {
		ErrorLogger.Printf("Cannot get %s because key value store is nil", key)
		return Entry{}, ErrKvStoreDoesNotExist
	}
	ds.Lock()
	defer ds.Unlock()
	existingData, ok := ds.kvStore[key]
	if !ok {
		InfoLogger.Printf("Cannot get %s because it does not exist", key)
		return Entry{}, ErrKeyNotFound
	} else {
		return NewEntry(key, existingData), nil
	}
}

func Authorized(data *Data, user string) bool {
	return data.owner == user
}
