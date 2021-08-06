// Package id 递增的ID，多节点请自行维护nodeID
package id

import (
	"basic/cipher"
	"basic/color"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"log"
	"strconv"
	"time"
)

var SId *server

var sNode *snowflake.Node

type server struct {
}

//String 获取string id
func (s server) String() string {
	return s.ToString(s.Int())
}

// Int 获取 int id
func (s server) Int() int64 {
	return sNode.Generate().Int64()
}

// ToString int转string
func (s server) ToString(id int64) string {
	return cipher.Base64EncryptInt64(id)
}

func (s server) ToInt(id string) int64 {
	return cipher.Base64DecryptBytesInt(id)
}

// Info id信息
func (s server) Info(id int64) (map[string]interface{}, error) {
	sId := snowflake.ParseInt64(id)
	tm := time.Unix(sId.Time()/1000, 0)
	return map[string]interface{}{
		"time": tm.Format("2006-01-02 15:04:05"),
	}, nil
}

// Decrypt 解码id
func (s server) Decrypt(str string) (map[string]interface{}, error) {
	id := cipher.Base64DecryptBytesInt(str)
	info, err := s.Info(id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	info["id"] = strconv.FormatInt(id, 10)
	return info, nil
}

// Encrypt 编码id
func (s server) Encrypt(idStr string) (interface{}, error) {
	int64Num, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return cipher.Base64EncryptInt64(int64Num), nil
}

type Server struct {
	Node int64
}

func (s Server) CreateNode() {
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
