package figure

import (
	"basic/id"
	"basic/oss"
	"basic/reptile/page"
	"testing"
)

func TestBaiduImage(t *testing.T) {
	page.Server{}.CreateServer()

	image, err := BaiduImage("苏州风景")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log(image)
}

func TestPrivateBaiduImage(t *testing.T) {
	//id
	id.Server{
		Node: 1,
	}.CreateNode()

	//云存储
	oss.Server{
		Endpoint:        "",
		AccessKeyId:     "",
		AccessKeySecret: "",
		BucketName:      "",
	}.CreateClient()
	//爬虫
	page.Server{}.CreateServer()

	image, err := PrivateBaiduImage("苏州风景")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log(image)
}
