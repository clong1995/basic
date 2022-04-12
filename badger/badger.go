package badger

import (
	"fmt"
	"github.com/clong1995/basic/color"
	badgerDB "github.com/dgraph-io/badger/v3"
	"log"
)

//https://gist.github.com/alexanderbez/d99fd0383ad57e991e9af9adcbb70b9d

type (
	Server struct {
	}

	server struct {
	}
)

var (
	Dadger *server
	bdb    *badgerDB.DB
)

func (s server) Get(namespace, key []byte) (value []byte, err error) {
	err = bdb.View(func(txn *badgerDB.Txn) error {
		var item *badgerDB.Item
		item, err = txn.Get(badgerNamespaceKey(namespace, key))
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (s server) Set(namespace, key, value []byte) (err error) {
	err = bdb.Update(func(txn *badgerDB.Txn) error {
		return txn.Set(badgerNamespaceKey(namespace, key), value)
	})
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (s server) Has(namespace, key []byte) (ok bool, err error) {
	_, err = s.Get(namespace, key)
	switch err {
	case badgerDB.ErrKeyNotFound:
		ok, err = false, nil
	case nil:
		ok, err = true, nil
	}
	return
}

func (s server) Del(namespace, key []byte) (err error) {
	err = bdb.Update(func(txn *badgerDB.Txn) error {
		return txn.Delete(badgerNamespaceKey(namespace, key))
	})
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (s Server) Run() {
	if Dadger != nil {
		return
	}
	opt := badgerDB.DefaultOptions("").WithInMemory(true)
	var dbErr error
	bdb, dbErr = badgerDB.Open(opt)
	if dbErr != nil {
		log.Fatal(color.Red, dbErr, color.Reset)
	}
	Dadger = new(server)
	color.Success("[badger] in memory connect success")
}

func badgerNamespaceKey(namespace, key []byte) []byte {
	prefix := []byte(fmt.Sprintf("%s/", namespace))
	return append(prefix, key...)
}
