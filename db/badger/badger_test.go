package badger

import (
	"testing"
)

func Test_badger(t *testing.T) {
	var namespace = []byte("namespace")
	var key = []byte("key")
	var value = []byte("value")
	Server{}.Run()
	//set
	err := Dadger.Set(namespace, key, value)
	if err != nil {
		t.Log(err)
		return
	}
	//get
	get, err := Dadger.Get(namespace, key)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(string(get))
	//has
	has, err := Dadger.Has(namespace, key)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(has)
	//delete
	err = Dadger.Del(namespace, key)
	if err != nil {
		t.Log(err)
		return
	}
	//has
	has, err = Dadger.Has(namespace, key)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(has)
}
