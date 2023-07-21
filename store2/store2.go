package store2

import "github.com/pkg/errors"

var (
	ErrNotFound = errors.New("not found")
)

const (
	bufferSize = 1 << 12
)

type actionAdd struct {
	key   string
	value string
}

type actionDel struct {
	key string
}

type actionGet struct {
	key   string
	value chan string
	err   chan error
}

type Storer struct {
	addChan chan *actionAdd
	delChan chan *actionDel
	getChan chan *actionGet
	data    map[string]string
}

func New() *Storer {
	storer := &Storer{
		addChan: make(chan *actionAdd, bufferSize),
		delChan: make(chan *actionDel, bufferSize),
		getChan: make(chan *actionGet, bufferSize),
		data:    make(map[string]string),
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
				action.err <- ErrNotFound
			} else {
				action.value <- value
			}
		}
	}
}

func (storer *Storer) Add(key, value string) {
	storer.addChan <- &actionAdd{key: key, value: value}
}

func (storer *Storer) Delete(key string) error {
	if _, err := storer.Get(key); err != nil {
		return errors.Wrap(err, "cannot delete")
	}

	storer.delChan <- &actionDel{key: key}

	return nil
}

func (storer *Storer) Get(key string) (string, error) {
	errChan := make(chan error, 1)
	valueChan := make(chan string, 1)

	defer close(errChan)
	defer close(valueChan)

	storer.getChan <- &actionGet{
		key:   key,
		err:   errChan,
		value: valueChan,
	}

	select {
	case err := <-errChan:
		return "", errors.Wrap(err, "cannot get")
	case value := <-valueChan:
		return value, nil
	}
}
