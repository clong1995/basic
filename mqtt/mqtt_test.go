package mqtt

import (
	"testing"
)

func TestServer_CreateClient(t *testing.T) {
	Server{
		Endpoint:        "onsmqtt.cn-beijing.aliyuncs.com",
		InstanceId:      "post-cn-i7m264v5w0g",
		AccessKeyId:     "",
		AccessKeySecret: "",
	}.Run()
	/*deviceCredential, err := Mqtt.Register("GID_dating", "dating-server")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("%+v", deviceCredential)*/

	/*_auth, err := Mqtt.AuthToken("GID_dating", "test_device", "dating-server", "R")
	if err != nil {
		return
	}
	t.Logf("%+v", _auth)*/
}
