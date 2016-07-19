package runnlib

import (
	"errors"
	"io"
)

type Store interface {
	Set(name string, r io.Reader) error
	Get(name string) (io.ReaderAt, int64, error)
}

type StoreFunc func(config interface{}) (Store, error)

var _stores map[string]StoreFunc

func AddStore(name string, fn StoreFunc) {
	_stores[name] = fn
}

func GetStore(name string, config interface{}) (Store, error) {
	if fn, ok := _stores[name]; ok {
		return fn(config)
	}
	return nil, errors.New("no store called: " + name)
}

func init() {
	_stores = make(map[string]StoreFunc)
}