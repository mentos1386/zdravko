package database

import (
	"time"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/pkg/errors"
)

type BadgerKeyValueStore struct {
	db *badger.DB
}

func NewBadgerKeyValueStore(path string) (*BadgerKeyValueStore, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open badger db")
	}
	return &BadgerKeyValueStore{db: db}, nil
}

func (b *BadgerKeyValueStore) Close() error {
	return b.db.Close()
}

func (b *BadgerKeyValueStore) Set(key string, value []byte, ttl time.Duration) error {
	return b.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), value).WithTTL(ttl)
		return txn.SetEntry(e)
	})
}

func (b *BadgerKeyValueStore) Increment(key string) (int, error) {
	var value int
	return value, b.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		value = int(valCopy[0]) + 1
		return txn.Set([]byte(key), []byte{byte(value)})
	})
}

func (b *BadgerKeyValueStore) Get(key string) ([]byte, error) {
	var value []byte
	return value, b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		valCopy, err := item.ValueCopy(value)
		if err != nil {
			return err
		}
		value = valCopy
		return nil
	})
}

func (b *BadgerKeyValueStore) Delete(key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (b *BadgerKeyValueStore) Keys(prefix string) ([]string, error) {
	var keys []string
	return keys, b.db.View(func(txn *badger.Txn) error {
		itr := txn.NewIterator(badger.DefaultIteratorOptions)
		defer itr.Close()
		for itr.Seek([]byte(prefix)); itr.ValidForPrefix([]byte(prefix)); itr.Next() {
			item := itr.Item()
			keys = append(keys, string(item.Key()))
		}
		return nil
	})
}
