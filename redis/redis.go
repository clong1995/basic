package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/clong1995/basic/color"
	"github.com/go-redis/redis/v9"
	"log"
	"time"
)

//https://studygolang.com/articles/32352
//https://www.tizi365.com/archives/290.html
//https://segmentfault.com/a/1190000021538684

//https://blog.csdn.net/ghosind/article/details/107327922
//https://pkg.go.dev/github.com/go-redis/redis/v8#section-readme
//https://pkg.go.dev/github.com/go-redis/redis#Client.Watch
//https://pkg.go.dev/github.com/go-redis/redis/v8

type (
	Server struct {
		Addr     string
		Password string
		DB       int
		Flush    bool
		Block    bool
	}
	server struct {
	}
)

var (
	Redis       *server
	redisClient *redis.Client
	//ctx         = context.Background()
)

//常用函数的封装

func (s server) Get(key string) (bytes []byte, err error) {
	return redisClient.Get(context.Background(), key).Bytes()
}

func (s server) Set(key string, value interface{}, expiration ...time.Duration) (err error) {
	exp := time.Duration(0)
	if len(expiration) == 1 {
		exp = expiration[0]
	}
	err = redisClient.Set(context.Background(), key, value, exp).Err()
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (s server) HGet(key, field string) (bytes []byte, err error) {
	if bytes, err = redisClient.HGet(context.Background(), key, field).Bytes(); err != nil {
		log.Println(err)
		return
	}
	return
}

func (s server) HSet(key, field string, value interface{}, expiration ...time.Duration) (err error) {
	cxt := context.Background()
	if err = redisClient.HSet(cxt, key, field, value).Err(); err != nil {
		log.Println(err)
		return
	}

	//过期时间
	if len(expiration) == 1 {
		exp := expiration[0]
		err = redisClient.Expire(cxt, key, exp).Err()
		if err != nil {
			log.Println(err)
			return
		}
	}

	return
}

func (s server) Exists(key string) (exists bool, err error) {
	i64, err := redisClient.Exists(context.Background(), key).Result()
	if err != nil {
		log.Println(err)
		return
	}
	return 1 == i64, nil
}

func (s server) Del(keys ...string) (count int64, err error) {
	return redisClient.Del(context.Background(), keys...).Result()
}

func (s server) HExists(key, field string) (exists bool, err error) {
	return redisClient.HExists(context.Background(), key, field).Result()
}

func (s server) HDel(key string, fields ...string) (count int64, err error) {
	return redisClient.HDel(context.Background(), key, fields...).Result()
}

func (s server) HSetStruct(key, field string, value interface{}, expiration ...time.Duration) (err error) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return
	}
	return s.HSet(key, field, jsonBytes, expiration...)
}

func (s server) HGetStruct(key, field string, v interface{}) (err error) {
	bytes, err := s.HGet(key, field)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

//删除过度封装
/*
func (s server) HMSet(key string, fields map[string]interface{}, expiration ...time.Duration) (err error) {
	//过期时间
	if len(expiration) == 1 {
		exp := expiration[0]
		err = redisClient.Expire(context.Background(), key, exp).Err()
		if err != nil {
			log.Println(err)
			return
		}
	}

	return redisClient.HMSet(context.Background(), key, fields).Err()
}

func (s server) HMGet(key string, fields ...string) (result []interface{}, err error) {
	return redisClient.HMGet(context.Background(), key, fields...).Result()
}

*/

// Client redis有丰富的api封装用起来不灵活，这里仅封装一些操作复杂的，简单的直接使用原生api
func (s server) Client() *redis.Client {
	return redisClient
}

func (s Server) Run() {
	if Redis != nil {
		return
	}

	if s.Addr == "" {
		s.Addr = ":6379"
	}

	//创建客户端
	redisClient = redis.NewClient(&redis.Options{
		Addr:     s.Addr,
		Password: s.Password, // no password set
		DB:       s.DB,       // use default DB
	})
	/*defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {
			log.Fatal(color.Red, err, color.Reset)
		}
	}(redisClient)*/

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(color.Red, err, color.Reset)
		return
	}

	//清空所有数据
	if s.Flush {
		_, err = redisClient.FlushDB(context.Background()).Result()
		if err != nil {
			log.Fatal(color.Red, err, color.Reset)
		}
	}
	Redis = new(server)
	color.Success(fmt.Sprintf("[redis] connect %s db %d success", s.Addr, s.DB))

	if s.Block {
		select {}
	}
}
