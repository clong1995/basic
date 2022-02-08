// Package oss 递增的ID，多节点请自行维护nodeID
package oss

import (
	"basic/color"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type (
	Server struct {
		Endpoint        string
		AccessKeyId     string
		AccessKeySecret string
		BucketName      string
	}
	UploadUrl struct {
		Url            string `json:"url"`
		OSSAccessKeyId string `json:"OSSAccessKeyId"`
		Policy         string `json:"policy"`
		Signature      string `json:"signature"`
		Key            string `json:"key"`
	}
	server struct {
		bucketName string
	}
)

var (
	Oss       *server
	ossClient *oss.Client
	bucket    *oss.Bucket
)

func (s server) PutSignPolicyFileIdURL(fId string) (url UploadUrl, err error) {
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
		log.Println(err)
		return
	}
	policyStr := base64.StdEncoding.EncodeToString(policyBytes)

	h := hmac.New(sha1.New, []byte(ossClient.Config.AccessKeySecret))
	_, err = io.WriteString(h, policyStr)
	if err != nil {
		log.Println(err)
		return
	}
	url.Url = "https://" + s.bucketName + "." + ossClient.Config.Endpoint
	url.OSSAccessKeyId = ossClient.Config.AccessKeyID
	url.Policy = policyStr
	url.Signature = base64.StdEncoding.EncodeToString(h.Sum(nil))
	url.Key = fId
	return
}

func (s server) PutSignFileIdURL(fId string) (url string, err error) {
	url, err = bucket.SignURL(fId, oss.HTTPPut, 60)
	if err != nil {
		log.Println(err)
		return
	}
	return
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
	scheme := "https"
	netLoc := ossClient.Config.Endpoint[len(scheme+"://"):]
	return fmt.Sprintf("%s://%s.%s/%s", scheme, bucket.BucketName, netLoc, name)
}

// UploadBase64 上传base64
func (s server) UploadBase64(path, value string) (string, error) {
	fileContentPosition := strings.Index(value, ",")
	uploadBaseString := value[fileContentPosition+1:]
	decodeString, err := base64.StdEncoding.DecodeString(uploadBaseString)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(decodeString)
	err = bucket.PutObject(path, buf)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return path, nil
}

// UploadUrl 上传url
func (s server) UploadUrl(path, url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer func(body io.ReadCloser) {
		err = body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)

	err = bucket.PutObject(path, res.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return path, nil
}

func (s Server) Run() {
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
	color.Success(fmt.Sprintf("[oss] open %s/%s handle success", s.Endpoint, s.BucketName))
}
