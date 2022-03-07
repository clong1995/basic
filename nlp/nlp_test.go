package nlp

import "testing"

func Test_server_Keywords(t *testing.T) {
	Server{
		AccessKeyId:     "",
		AccessKeySecret: "",
	}.Run()
	result := NLP.Keywords("北京市两天三地游")

	t.Logf("%+v", result)
}
