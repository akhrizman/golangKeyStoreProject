package datasource

import (
	"errors"
	"httpstore/log4g"
	"sync"
	"time"
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
	depth   int
	mutex   *sync.RWMutex
}

func NewDatasource(depth int) Datasource {
	log4g.Info.Println("Created New Datasource with keystore")
	kvStore := make(map[Key]Data)
	mutex := sync.RWMutex{}
	return Datasource{kvStore, depth, &mutex}
}

// EvictLru Deletes the Least Recently Used (LRU) key if keystore size is at depth
func (ds *Datasource) EvictLru() {
	if ds.Size() == ds.depth {
		delete(ds.kvStore, ds.getLruKey())
	}
}

// getLruKey Retrieves the key with the oldest timestamp from the keystore
func (ds *Datasource) getLruKey() Key {
	var lruKey Key
	var lruDate time.Time
	isFirstKeyChecked := true
	for key, data := range ds.kvStore {
		// Basic logic: Compare each key to the previous key.  First key checked is set automatically
		if isFirstKeyChecked {
			lruKey = key
			lruDate = data.lastUsed
			isFirstKeyChecked = false
		} else if data.lastUsed.Before(lruDate) {
			lruKey = key
			lruDate = data.lastUsed
		}
	}
	return lruKey
}

func (ds *Datasource) SetDepth(depth int) {
	ds.depth = depth
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
		log4g.Error.Println("Did not get actual size because key value store is nil")
		return 0
	}
	return len(ds.kvStore)
}

func (ds *Datasource) EmptyKvStore() error {
	if ds.isClosed() {
		log4g.Error.Println("Cannot make key value store nil when already nil")
		return ErrKvStoreDoesNotExist
	} else {
		log4g.Info.Println("Removed key store from datasource")
		ds.kvStore = nil
	}
	return nil
}

func (ds *Datasource) CreateKvStore() error {
	if ds.isOpen() {
		log4g.Error.Println("Cannot instantiate key value store when already instantiated")
		return ErrKvStoreAlreadyExists
	} else {
		log4g.Info.Println("Created new key store for datasource")
		ds.kvStore = map[Key]Data{}
	}
	return nil
}

func (ds *Datasource) Put(key Key, newData Data) error {
	if ds.isClosed() {
		log4g.Error.Printf("Cannot insert %s because key value store is nil", key)
		return ErrKvStoreDoesNotExist
	}
	ds.Lock()
	defer ds.Unlock()
	existingData, ok := ds.kvStore[key]
	if ok && !Authorized(&existingData, newData.owner) {
		log4g.Info.Printf("Cannot update %s because owners do not match", key)
		return ErrValueUpdateForbidden
	} else {
		existingData.SetToCurrentTime()
		ds.EvictLru()
		ds.kvStore[key] = newData
	}
	return nil
}

func (ds *Datasource) Contains(key Key) (bool, error) {
	if ds.isClosed() {
		log4g.Error.Printf("Cannot check if %s exists because key value store is nil", key)
		return false, ErrKvStoreDoesNotExist
	}
	ds.RLock()
	defer ds.RUnlock()
	_, ok := ds.kvStore[key]
	return ok, nil
}

func (ds *Datasource) Get(key Key) (*Data, error) {
	if ds.isClosed() {
		log4g.Error.Printf("Cannot get %s because key value store is nil", key)
		return nil, ErrKvStoreDoesNotExist
	}
	ds.Lock()
	defer ds.Unlock()
	existingData, ok := ds.kvStore[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	existingData.SetToCurrentTime()
	ds.kvStore[key] = existingData
	return &existingData, nil
}

func (ds *Datasource) Delete(key Key, user string) error {
	if ds.isClosed() {
		log4g.Error.Printf("Cannot delete %s because key value store is nil", key)
		return ErrKvStoreDoesNotExist
	}
	ds.Lock()
	defer ds.Unlock()
	existingData, ok := ds.kvStore[key]
	if !ok {
		log4g.Error.Printf("Cannot delete %s because it does not exist", key)
		return ErrKeyNotFound
	} else if Authorized(&existingData, user) {
		delete(ds.kvStore, key)
		log4g.Info.Printf("Deleted %s with value <%s>", key, existingData.value)
		return nil
	} else {
		log4g.Info.Printf("Cannot update %s because owners do not match", key)
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
		log4g.Error.Printf("Cannot get %s because key value store is nil", key)
		return Entry{}, ErrKvStoreDoesNotExist
	}
	ds.Lock()
	defer ds.Unlock()
	existingData, ok := ds.kvStore[key]
	if !ok {
		return Entry{}, ErrKeyNotFound
	} else {
		return NewEntry(key, existingData), nil
	}
}

func Authorized(data *Data, user string) bool {
	// At this point, if user is admin, then any PUT/DELETE should
	// behave as if being performed by the original key's owner
	if user == "admin" {
		user = data.owner
	}
	return data.owner == user
}
