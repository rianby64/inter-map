package store1

import (
	"sync"

	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type Storer struct {
	locker *sync.RWMutex
	data   map[string]string
}

func New() *Storer {
	storer := &Storer{
		locker: &sync.RWMutex{},
		data:   make(map[string]string),
	}

	return storer
}

func (storer *Storer) Add(key, value string) {
	storer.locker.Lock()

	defer storer.locker.Unlock()

	storer.data[key] = value
}

func (storer *Storer) Delete(key string) error {
	storer.locker.Lock()

	defer storer.locker.Unlock()

	if _, hasKey := storer.data[key]; !hasKey {
		return ErrNotFound
	}

	delete(storer.data, key)

	return nil
}

func (storer *Storer) Get(key string) (string, error) {
	storer.locker.RLock()

	defer storer.locker.RUnlock()

	value, hasKey := storer.data[key]
	if !hasKey {
		return "", ErrNotFound
	}

	return value, nil
}
