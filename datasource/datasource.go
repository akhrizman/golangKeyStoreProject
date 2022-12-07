package datasource

import (
	"errors"
	"fmt"
	"sync"
)

// Key type
type Key string

func (key Key) String() string {
	return fmt.Sprintf("Key(%s)", string(key))
}

// Data to contain owner and value
type Data struct {
	owner string `json:"owner"`
	value string `json:"value"`
}

func NewData(owner string, value string) Data {
	return Data{owner: owner, value: value}
}
func (d *Data) GetOwner() string {
	return d.owner
}
func (d *Data) SetOwner(owner string) {
	d.owner = owner
}
func (d *Data) GetValue() string {
	return d.value
}
func (d *Data) SetValue(value string) {
	d.value = value
}

// Errors
var (
	ErrKeyNotFound          = errors.New("key not found")
	ErrKvStoreAlreadyExists = errors.New("key value store has already been initialized")
	ErrKvStoreDoesNotExist  = errors.New("key value store has not been initialized")
	ErrValueUpdateForbidden = errors.New("value update not allowed, wrong owner")
)

// Datasource and corresponding methods
type Datasource struct {
	kvStore map[Key]Data
	mutex   sync.RWMutex
}

func NewDatasource() Datasource {
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

func (ds *Datasource) EmptyKvStore() error {
	if ds.isClosed() {
		return ErrKvStoreDoesNotExist
	} else {
		ds.kvStore = nil
	}
	return nil
}

func (ds *Datasource) CreateKvStore() error {
	if ds.isOpen() {
		return ErrKvStoreAlreadyExists
	} else {
		ds.kvStore = map[Key]Data{}
	}
	return nil
}

func (ds *Datasource) Size() int {
	return len(ds.kvStore)
}

func (ds *Datasource) Put(key Key, data Data) error {
	if ds.isClosed() {
		return ErrKvStoreDoesNotExist
	}
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	existingData, ok := ds.kvStore[key]
	if ok && existingData.owner != data.owner {
		return ErrValueUpdateForbidden
	} else {
		ds.kvStore[key] = data
	}
	return nil
}

func (ds *Datasource) Contains(key Key) (bool, error) {
	if ds.isClosed() {
		return false, ErrKvStoreDoesNotExist
	}
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	_, ok := ds.kvStore[key]
	return ok, nil
}

func (ds *Datasource) Get(key Key) (*Data, error) {
	if ds.isClosed() {
		return nil, ErrKvStoreDoesNotExist
	}
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	existingData, ok := ds.kvStore[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return &existingData, nil
}

func (ds *Datasource) Delete(key Key) error {
	if ds.isClosed() {
		return ErrKvStoreDoesNotExist
	}
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	_, ok := ds.kvStore[key]
	if !ok {
		return ErrKeyNotFound
	} else {
		delete(ds.kvStore, key)
		return nil
	}
}
