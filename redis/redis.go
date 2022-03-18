package redis

import (
	"basic/color"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

//https://studygolang.com/articles/32352
//https://www.tizi365.com/archives/290.html
//https://segmentfault.com/a/1190000021538684

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
	ctx         = context.Background()
)

//

func (s server) Set(key string, value interface{}, expiration ...time.Duration) (err error) {
	exp := time.Duration(0)
	if len(expiration) == 1 {
		exp = expiration[0]
	}
	err = redisClient.Set(ctx, key, value, exp).Err()
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (s server) Get(key string) (bytes []byte, err error) {
	return redisClient.Get(ctx, key).Bytes()
}

func (s server) Exists(key string) (exists bool, err error) {
	i64, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		log.Println(err)
		return
	}
	return 1 == i64, nil
}

func (s server) Del(keys ...string) (count int64, err error) {
	count, err = redisClient.Del(ctx, keys...).Result()
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (s server) HSetStruct(key, field string, value interface{}, expiration ...time.Duration) (err error) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return
	}
	return s.HSet(key, field, jsonBytes, expiration...)
}

func (s server) HSet(key, field string, value interface{}, expiration ...time.Duration) (err error) {
	//过期时间
	if len(expiration) == 1 {
		exp := expiration[0]
		err = redisClient.Expire(ctx, key, exp).Err()
		if err != nil {
			log.Println(err)
			return
		}
	}
	if err = redisClient.HSet(ctx, key, field, value).Err(); err != nil {
		log.Println(err)
		return
	}
	return
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

func (s server) HGet(key, field string) (bytes []byte, err error) {
	if bytes, err = redisClient.HGet(ctx, key, field).Bytes(); err != nil {
		log.Println(err)
		return
	}
	return
}

func (s server) HExists(key, field string) (exists bool, err error) {
	exists, err = redisClient.HExists(ctx, key, field).Result()
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (s server) HDel(key string, fields ...string) (count int64, err error) {
	count, err = redisClient.HDel(ctx, key, fields...).Result()
	if err != nil {
		log.Println(err)
		return
	}
	return
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
	//清空所有数据
	if s.Flush {
		_, err := redisClient.FlushDB(ctx).Result()
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
