// Package oss 递增的ID，多节点请自行维护nodeID
package oss

import (
	"basic/color"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
	"time"
)

var Oss *server
var ossClient *oss.Client
var bucket *oss.Bucket

type server struct {
	bucketName string
}

/*// PutSignURL 签名直传
func (s server) PutSignURL() (fileId string, url string, err error) {
	fileId = id.SId.String()
	//url, err = bucket.SignURL(fileId, oss.HTTPPut, 60)
	url, err = bucket.SignURL(fileId, oss.HTTPPut, 60)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
// PutSignFileIdURL 签名直传
func (s server) PutSignFileIdURL(fId string) (url string, err error) {
	url, err = bucket.SignURL(fId, oss.HTTPPut, 60)
	if err != nil {
		log.Println(err)
		return
	}
	return
}*/

func (s server) PutSignFileIdURL(fId string) (param map[string]string, err error) {
	policy := map[string]interface{}{
		"expiration": time.Now().Add(60 * time.Second).Format("2006-01-02T15:04:05.999Z"),
		"conditions": []interface{}{
			map[string]string{
				"bucket": s.bucketName,
			},
			[]string{"eq", "$key", fId},
		},
	}
	policyBytes, err := json.Marshal(&policy)
	if err != nil {
		return
	}
	policyStr := base64.StdEncoding.EncodeToString(policyBytes)

	h := hmac.New(sha1.New, []byte(ossClient.Config.AccessKeySecret))
	_, err = io.WriteString(h, policyStr)
	if err != nil {
		return
	}
	return map[string]string{
		"url":            "https://" + s.bucketName + "." + ossClient.Config.Endpoint,
		"OSSAccessKeyId": ossClient.Config.AccessKeyID,
		"policy":         policyStr,
		"signature":      base64.StdEncoding.EncodeToString(h.Sum(nil)),
		"key":            fId,
	}, nil
}

// GetSignURL TODO 未测试，很大可能有问题。
// 使用签名URL将OSS文件下载到流。
/*func (s server) GetSignURL(name string, expiredInSec ...int64) (string, error) {
	if name == "" {
		return "", fmt.Errorf("name不得为空\n")
	}
	var sec int64 = 60
	if len(expiredInSec) == 1 && expiredInSec[0] > sec {
		sec = expiredInSec[0]
	}

	return bucket.SignURL(name, oss.HTTPGet, sec)
}*/

// GetURL 将OSS文件下载到流。
func (s server) GetURL(name string) string {
	if name == "" {
		return ""
	}
	return fmt.Sprintf("https://%s.%s/%s", bucket.BucketName, ossClient.Config.Endpoint, name)
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

	//创建对象，禁止项目访问这些信息
	Oss = &server{
		bucketName: s.BucketName,
	}
	color.Success(fmt.Sprintf("[oss] open %s %s handle success", s.Endpoint, s.BucketName))
}
