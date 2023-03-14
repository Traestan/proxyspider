package storage

import (
	badger "github.com/dgraph-io/badger/v3"
)

// BadgerStorage структура для хранения в файле
type BadgerStorage struct {
	*badger.DB
}

// NewBadgerStorage решение для записи в badgerdb
func NewBadgerStorage(path string) *BadgerStorage {
	db, err := badger.Open(badger.DefaultOptions("./badger"))
	if err != nil {
		panic(err)
	}
	return &BadgerStorage{db}
}
