package redis

import (
	"testing"
)

func Test_server_HSet(t *testing.T) {
	Server{
		Addr:     "192.168.1.100:6379",
		Password: "Hzm2022YCL..",
		DB:       0,
	}.Run()

	err := Redis.HSet("key1", "field1", "value1")
	if err != nil {
		t.Log(err)
		return
	}
}
