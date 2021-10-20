package page

import (
	"testing"
)

func TestDoc(t *testing.T) {
	Server{}.CreateServer()
	doc, err := Page.Doc("https://image.baidu.com/search/index?tn=baiduimage&word=苏州风景&copyright=0", "#imgid > div:nth-child(1) > ul")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log(doc)
}
