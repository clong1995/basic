package oss

import "testing"

func Test_server_GetURL(t *testing.T) {
	Server{
		Endpoint:        "oss-cn-beijing",
		AccessKeyId:     "",
		AccessKeySecret: "",
		BucketName:      "project-lovers",
	}.CreateClient()
	url := Oss.GetURL("/dir/test.jpg")
	t.Log(url)
}
