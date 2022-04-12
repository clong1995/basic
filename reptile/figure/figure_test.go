package figure

import (
	"github.com/clong1995/basic/id"
	"github.com/clong1995/basic/oss"
	"github.com/clong1995/basic/reptile/page"
	"testing"
)

func TestBaiduImage(t *testing.T) {
	page.Server{}.Run()

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
	}.Run()

	//云存储
	oss.Server{
		Endpoint:        "",
		AccessKeyId:     "",
		AccessKeySecret: "",
		BucketName:      "",
	}.Run()
	//爬虫
	page.Server{}.Run()

	image, err := PrivateBaiduImage("", "苏州风景")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log(image)
}
