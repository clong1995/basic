// Package oss 递增的ID，多节点请自行维护nodeID
package oss

import (
	"basic/color"
	"basic/id"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
)

var Oss *server
var ossClient *oss.Client
var bucket *oss.Bucket

type server struct {
}

// PutSignURL 签名直传
func (s server) PutSignURL(objectType string) (fileId string, url string, err error) {
	fileId = id.SId.String()
	url, err = bucket.SignURL(fileId, oss.HTTPPut, 60)
	if err == nil {
		return
	}
	return
}

// GetSignURL 使用签名URL将OSS文件下载到流。
func (s server) GetSignURL(name string, expiredInSec ...int64) (string, error) {
	if name == "" {
		return "", fmt.Errorf("name不得为空\n")
	}
	var sec int64 = 60
	if len(expiredInSec) == 1 && expiredInSec[0] > sec {
		sec = expiredInSec[0]
	}

	return bucket.SignURL(name, oss.HTTPGet, sec)
}

// GetURL 使用签名URL将OSS文件下载到流。
func (s server) GetURL(name string) string {
	if name == "" {
		return ""
	}
	return fmt.Sprintf("https://%s.%s.aliyuncs.com/%s", bucket.BucketName, ossClient.Config.Endpoint, name)
}

type Server struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
}

func (s Server) CreateClient() {
	if Oss != nil {
		return
	}
	var err error
	//创建客户端
	ossClient, err = oss.New(s.Endpoint, s.AccessKeyId, s.AccessKeySecret)
	if err != nil {
		log.Fatalln(color.Red, err, color.Reset)
	}

	// 获取存储空间。
	bucket, err = ossClient.Bucket(s.BucketName)
	if err != nil {
		log.Fatalln(color.Red, err, color.Reset)
	}

	//创建对象
	Oss = new(server)
	color.Success(fmt.Sprintf("[oss] open %s %s handle success", s.Endpoint, s.BucketName))
}
