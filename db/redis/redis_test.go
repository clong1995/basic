package redis

import "testing"

func Test_server_SetStruct(t *testing.T) {
	Server{
		Addr:     "",
		Password: "",
		DB:       0,
	}.Run()

	type data struct {
		Name        string `json:"name"`
		Age         int    `json:"age"`
		EnglishName string `json:"english_name"`
	}
	err := Redis.SetStruct("testData", "user", data{
		Name:        "于成龙",
		Age:         25,
		EnglishName: "eleven",
	})
	if err != nil {
		t.Log(err)
	}
}

func Test_server_GetStruct(t *testing.T) {
	Server{
		Addr:     "redis-13579.c278.us-east-1-4.ec2.cloud.redislabs.com:13579",
		Password: "clong11429ycl.YU",
		DB:       0,
	}.Run()

	type data struct {
		Name        string `json:"name"`
		Age         int    `json:"age"`
		EnglishName string `json:"english_name"`
	}
	d := new(data)
	err := Redis.GetStruct("testData", "user", d)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(d)
}

func Test_server_ExistsStruct(t *testing.T) {
	Server{
		Addr:     "redis-13579.c278.us-east-1-4.ec2.cloud.redislabs.com:13579",
		Password: "clong11429ycl.YU",
		DB:       0,
	}.Run()
	exists, err := Redis.ExistsStruct("testData", "user")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(exists)
}

func Test_server_Del(t *testing.T) {
	Server{
		Addr:     "redis-13579.c278.us-east-1-4.ec2.cloud.redislabs.com:13579",
		Password: "clong11429ycl.YU",
		DB:       0,
	}.Run()
	count, err := Redis.Del("testData", "user")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(count)
}
