package store3

import "github.com/pkg/errors"

var (
	ErrNotFound = errors.New("not found")
)

const (
	bufferSize = 1 << 4
)

type actionAdd struct {
	key   string
	value string
}

type actionDel struct {
	key string
}

type actionGet struct {
	key string
}

type Storer struct {
	addChan chan actionAdd
	delChan chan actionDel
	getChan chan actionGet

	getResultValueChan chan string
	getResultErrChan   chan error

	data map[string]string
}

func New() *Storer {
	storer := &Storer{
		addChan: make(chan actionAdd, bufferSize),
		delChan: make(chan actionDel, bufferSize),
		getChan: make(chan actionGet, bufferSize),
		data:    make(map[string]string),

		getResultValueChan: make(chan string, 1),
		getResultErrChan:   make(chan error, 1),
	}

	go storer.processActions()

	return storer
}

// https://go-proverbs.github.io/
// Don't communicate by sharing memory, share memory by communicating.
func (storer *Storer) processActions() {
	for {
		select {
		case action := <-storer.addChan:
			storer.data[action.key] = action.value
		case action := <-storer.delChan:
			delete(storer.data, action.key)
		case action := <-storer.getChan:
			value, hasKey := storer.data[action.key]
			if !hasKey {
				storer.getResultErrChan <- ErrNotFound
			} else {
				storer.getResultValueChan <- value
			}
		}
	}
}

func (storer *Storer) Add(key, value string) {
	storer.addChan <- actionAdd{key: key, value: value}
}

func (storer *Storer) Delete(key string) error {
	if _, err := storer.Get(key); err != nil {
		return errors.Wrap(err, "cannot delete")
	}

	storer.delChan <- actionDel{key: key}

	return nil
}

func (storer *Storer) Get(key string) (string, error) {
	storer.getChan <- actionGet{
		key: key,
	}

	select {
	case err := <-storer.getResultErrChan:
		return "", errors.Wrap(err, "cannot get")
	case value := <-storer.getResultValueChan:
		return value, nil
	}
}
