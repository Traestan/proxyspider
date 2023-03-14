package storage

import (
	"fmt"
	"time"

	badger "github.com/dgraph-io/badger/v3"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"gitlab.com/likekimo/goproxyspider/pkg/storage"
	"go.uber.org/zap"
)

//boltStorage структура для хранения в файле
type badgerStorage struct {
	db     *storage.BadgerStorage
	logger *logger.Logger
	count  int
}

//BoltStorage решение для записи в boltdb
func BadgerStorage(logger *logger.Logger, path string) Storage {
	bdb := storage.NewBadgerStorage(path)

	svc := &badgerStorage{
		logger: logger,
		db:     bdb,
	}
	//beforeStat :=
	svc.count = svc.getstat()

	logger.Info("Badger init")

	return svc
}
func (bg badgerStorage) CheckStorage() error {
	err := bg.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte("1"), []byte(""))
		err := txn.SetEntry(e)
		return err
	})
	if err != nil {
		bg.logger.Error("check storage", zap.Error(fmt.Errorf("could not create root bucket: %v", err)))
		return err
	}

	return nil
}

func (bg badgerStorage) WriteStorage(source string) error {
	err := bg.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(source), []byte(time.Now().Format(time.RFC3339)))
		return err
	})
	return err
}
func (bg badgerStorage) Stat() error {
	stat := bg.getstat()
	bg.logger.Info("Count badger Storage", zap.Int("after", stat), zap.Int("before", bg.count))

	return nil
}

func (bg badgerStorage) getstat() int {
	bg.logger.Debug("Count badger Storage")
	var countBefore int
	err := bg.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			countBefore++
		}
		return nil
	})
	if err != nil {
		bg.logger.Debug("Count badger Storage", zap.Error(fmt.Errorf("stat not vwork: %v", err)))
		return countBefore
	}

	return countBefore
}
