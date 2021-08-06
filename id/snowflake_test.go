package id

import (
	"testing"
)

func TestInfo(t *testing.T) {
	Server{
		Node: 1,
	}.CreateNode()

	for i := 0; i < 100; i++ {
		//生成字符串id
		s := SId.String()
		t.Log(s)
		//转整型
		id := SId.ToInt(s)
		t.Log(id)
	}
}
