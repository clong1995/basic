// Package figure 根据文字获取一张图片
package figure

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/clong1995/basic/oss"
	"github.com/clong1995/basic/reptile/page"
	"log"
	"strings"
)

type Result struct {
	Type  string
	Value string
}

func BaiduImage(name string) (res Result, err error) {
	doc, err := page.Page.Doc(fmt.Sprintf("https://image.baidu.com/search/index?tn=baiduimage&word=%s&copyright=0", name), "#imgid > div:nth-child(1) > ul")
	if err != nil {
		log.Println(err)
		return
	}

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(doc))
	if err != nil {
		log.Println(err)
		return
	}

	selection := dom.Find(".imgitem")

	if selection.Length() < 0 {
		err = fmt.Errorf("无图片")
		return
	}

	first := selection.Eq(1)

	imgTag := first.Find("img")
	val, exists := imgTag.Attr("src")
	if !exists {
		return
	}

	res.Type = "url"
	if strings.HasPrefix(val, "data:image/jpeg;base64,") {
		res.Type = "base64"
	}
	res.Value = val
	return
}

func PrivateBaiduImage(path, name string) (url string, err error) {
	image, err := BaiduImage(name)
	if err != nil {
		log.Println(err)
		return
	}
	if image.Type == "base64" {
		//上传base64
		url, err = oss.Oss.UploadBase64(path, image.Value)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		//上传url
		url, err = oss.Oss.UploadUrl(path, image.Value)
		if err != nil {
			log.Println(err)
			return
		}
	}

	return
}
