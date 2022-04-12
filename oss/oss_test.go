package oss

import (
	"github.com/clong1995/basic/id"
	"testing"
)

func Test_server_GetURL(t *testing.T) {
	Server{
		Endpoint:        "",
		AccessKeyId:     "",
		AccessKeySecret: "",
		BucketName:      "",
	}.Run()
	url := Oss.GetURL("/dir/test.jpg")
	t.Log(url)
}

func TestUploadBase64(t *testing.T) {
	//id
	id.Server{
		Node: 1,
	}.Run()
	//
	Server{
		Endpoint:        "",
		AccessKeyId:     "",
		AccessKeySecret: "",
		BucketName:      "",
	}.Run()
	//
	url, err := Oss.UploadBase64("", "data:image/jpeg;base64,xxxxx")
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(url)
}
