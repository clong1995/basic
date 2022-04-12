// Package id 递增的ID，多节点请自行维护nodeID
package id

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/clong1995/basic/color"
	"log"
	"time"
)

var (
	SId   *server
	sNode *snowflake.Node
)

type (
	Server struct {
		Node int64
	}
	server struct {
	}
)

//String 获取string id
func (s server) String() string {
	return sNode.Generate().Base58()
}

// Int 获取 int id
func (s server) Int() int64 {
	return sNode.Generate().Int64()
}

// ToString int转string
func (s server) ToString(id int64) string {
	sId := snowflake.ParseInt64(id)
	return sId.Base58()
}

func (s server) ToInt(id string) int64 {
	sID, err := snowflake.ParseBase58([]byte(id))
	if err != nil {
		log.Println(err)
		return 0
	}

	return sID.Int64()
}

// Test 测试函数
func (s server) Test() {
	//base58
	b58 := s.String()
	log.Println(b58)

	//base58 to int64
	i64 := s.ToInt(b58)
	log.Println(i64)

	//int64 to base58
	b58 = s.ToString(i64)
	log.Println(b58)

	log.Println("================")

	//int64
	i64 = s.Int()
	log.Println(i64)
	//int64 to base58
	b58 = s.ToString(i64)
	log.Println(b58)
	//base58 to int64
	i64 = s.ToInt(b58)
	log.Println(i64)
}

// Info id信息
// DEPRECATED: the below function will be removed in a future release.
func (s server) Info(id int64) (map[string]interface{}, error) {
	sId := snowflake.ParseInt64(id)
	tm := time.Unix(sId.Time()/1000, 0)
	return map[string]interface{}{
		"time": tm.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s Server) Run() {
	//防止多次创建
	if SId != nil {
		return
	}
	var err error
	sNode, err = snowflake.NewNode(s.Node)
	//id生成器创建失败，直接退出
	if err != nil {
		log.Fatalln(color.Red, err, color.Reset)
	}
	//创建对象
	SId = new(server)
	color.Success(fmt.Sprintf("[snowflake] create node %d success", s.Node))
}
