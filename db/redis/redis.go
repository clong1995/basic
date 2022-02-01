package redis

import (
	"basic/color"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

//https://studygolang.com/articles/32352
//https://www.tizi365.com/archives/290.html

type (
	Server struct {
		Addr     string
		Password string
		DB       int
		Flush    bool
	}
	server struct {
	}
)

var (
	Redis       *server
	redisClient *redis.Client
	ctx         = context.Background()
)

func (s server) SetStruct(key, field string, value interface{}) (err error) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return
	}
	err = redisClient.HSet(ctx, key, field, string(jsonBytes)).Err()
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (s server) GetStruct(key, field string, v interface{}) (err error) {
	bytes, err := redisClient.HGet(ctx, key, field).Bytes()
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

func (s server) ExistsStruct(key, field string) (exists bool, err error) {
	exists, err = redisClient.HExists(ctx, key, field).Result()
	if err != nil {
		log.Println(err)
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (s server) Del(key string, keys ...string) (count int64, err error) {
	count, err = redisClient.HDel(ctx, key, keys...).Result()
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
}
