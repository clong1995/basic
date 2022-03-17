package id

import (
	"testing"
)

func TestToTest(t *testing.T) {
	Server{
		Node: 1,
	}.Run()

	/*for i := 0; i < 100; i++ {
		//生成字符串id
		s := SId.String()
		t.Log(s)
		//转整型
		id := SId.ToInt(s)
		t.Log(id)
	}*/

	strId := SId.ToString(1503985520630304770)
	t.Log(strId)
	//SId.Test()
}
