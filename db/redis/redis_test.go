package redis

import "testing"

func Test_server_Set(t *testing.T) {
	Server{
		Addr:     "",
		Password: "",
		DB:       0,
	}.Run()

	err := Redis.Set("test", "user")
	if err != nil {
		t.Log(err)
	}
}

func Test_server_Get(t *testing.T) {
	Server{
		Addr:     "",
		Password: "",
		DB:       0,
	}.Run()

	bytes, err := Redis.Get("test")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(bytes)
}

func Test_server_Exists(t *testing.T) {
	Server{
		Addr:     "",
		Password: "",
		DB:       0,
	}.Run()
	exists, err := Redis.Exists("test")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(exists)
}

func Test_server_Del(t *testing.T) {
	Server{
		Addr:     "",
		Password: "",
		DB:       0,
	}.Run()
	count, err := Redis.Del("test")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(count)
}

func Test_server_HSetStruct(t *testing.T) {
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
	err := Redis.HSetStruct("testData", "user", data{
		Name:        "于成龙",
		Age:         25,
		EnglishName: "eleven",
	})
	if err != nil {
		t.Log(err)
	}
}

func Test_server_HGetStruct(t *testing.T) {
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
	d := new(data)
	err := Redis.HGetStruct("testData", "user", d)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(d)
}

func Test_server_HExists(t *testing.T) {
	Server{
		Addr:     "",
		Password: "",
		DB:       0,
	}.Run()
	exists, err := Redis.HExists("testData", "user")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(exists)
}

func Test_server_HDel(t *testing.T) {
	Server{
		Addr:     "",
		Password: "",
		DB:       0,
	}.Run()
	count, err := Redis.HDel("testData", "user")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(count)
}
