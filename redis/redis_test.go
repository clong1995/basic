package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
	"testing"
)

func Test_server_ZRangeByScore(t *testing.T) {
	Server{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}.Run()

	data := make([]*redis.Z, 0)

	data = append(data, &redis.Z{Score: 0, Member: "tizi0"})
	data = append(data, &redis.Z{Score: 5, Member: "tizi5"})
	data = append(data, &redis.Z{Score: 10, Member: "tizi10"})
	data = append(data, &redis.Z{Score: 15, Member: "tizi15"})
	data = append(data, &redis.Z{Score: 20, Member: "tizi20"})

	err := Redis.Client().ZAdd(context.Background(), "zaddtest", data...).Err()
	if err != nil {
		t.Log(err)
	}

	op := &redis.ZRangeBy{
		Min:    "0",    // 最小分数
		Max:    "+inf", // 最大分数
		Offset: 0,      // 类似sql的limit, 表示开始偏移量
		Count:  23,     // 一次返回多少数据
	}

	result, err := Redis.Client().ZRangeByScore(context.Background(), "zaddtest", op).Result()
	if err != nil {
		t.Log(err)
	}

	for _, val := range result {
		t.Log(val)
	}
}

//https://geektutu.com/post/hpg-concurrency-control.html控制并发量
//多线程并发修改
func Test_server_WatchMutil(t *testing.T) {
	Server{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}.Run()
	//初始化key
	key := "mutil"
	err := Redis.Set(key, 0)
	if err != nil {
		t.Error(err)
		return
	}

	//多协程
	const routineCount = 5
	// 干活干
	/*worker := func() {
		defer wg.Done()

		// 干活干活干活
		var value []byte
		var cacheErr error

		value, cacheErr = Redis.Get(key)
		if cacheErr != nil {
			t.Error(cacheErr)
			return
		}
		var i64 int64
		i64, cacheErr = strconv.ParseInt(string(value), 10, 64)
		if cacheErr != nil {
			t.Error(cacheErr)
			return
		}
		i64 += 1
		cacheErr = Redis.Set(key, i64)
		if cacheErr != nil {
			t.Error(cacheErr)
			return
		}
	}*/

	// 累加
	increment := func(key string) (err error) {
		txf := func(tx *redis.Tx) (err error) {
			ctx := context.Background()
			//读取
			value, err := tx.Get(ctx, key).Int64()
			if err != nil {
				t.Error(err)
				return
			}
			//修改值
			value++

			// runs only if the watched keys remain unchanged
			_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) (err error) {
				//插入
				t.Log(value)
				err = pipe.Set(ctx, key, value, 0).Err()
				if err != nil {
					t.Error(err)
				}
				return
			})

			if err != nil {
				t.Error(err)
				return
			}

			return
		}

		err = Redis.Client().Watch(context.Background(), txf, key)
		if err != nil {
			t.Error(err)
			return
		}
		//TODO retries

		return
	}

	//并发修改
	wg := sync.WaitGroup{}
	wg.Add(routineCount)
	for i := 0; i < routineCount; i++ {
		//睡一会的结果，的确累加了，不睡补发得到正确累加
		//time.Sleep(1 * time.Second)
		go func() {
			defer wg.Done()
			if incrErr := increment(key); incrErr != nil {
				t.Error(incrErr)
			}
		}()
	}
	wg.Wait()
}
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

/*func Test_server_SetStruct(t *testing.T) {
	Server{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}.Run()

	type data struct {
		Name        string `json:"name"`
		Age         int    `json:"age"`
		EnglishName string `json:"english_name"`
	}

	err := Redis.Set("test_set_struct", data{
		Name:        "于成龙",
		Age:         25,
		EnglishName: "eleven",
	})
	if err != nil {
		t.Log(err)
	}
}*/

func Test_server_Get(t *testing.T) {
	Server{
		Addr:     "127.0.0.1:6379",
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
func Test_server_HGet(t *testing.T) {
	Server{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}.Run()

	bytes, err := Redis.HGet("test", "xxx")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(bytes)
}

/*func Test_server_HMGet(t *testing.T) {
	Server{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}.Run()

	result, err := Redis.HMGet("share", `{"datetime":"2022032215","offset":0}`, `bbb`)
	if err != nil {
		t.Log(err)
		return
	}
	for _, res := range result {
		t.Log(res)
	}
}*/

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

/*func Test_server_HSetStruct(t *testing.T) {
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
}*/

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
