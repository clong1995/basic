package request

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func HttpPostJson(url string, data interface{}, headers ...map[string]string) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//设置头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if len(headers) > 0 {
		for s, s2 := range headers[0] {
			req.Header.Set(s, s2)
		}
	}

	return doRequest(req)
}

func HttpGet(urlStr string, data map[string]string, headers ...map[string]string) ([]byte, error) {
	//url
	Url, err := url.Parse(urlStr)
	if err != nil {
		panic(err.Error())

	}

	//参数
	if data != nil {
		params := url.Values{}
		for s, s2 := range data {
			params.Set(s, s2)
		}
		Url.RawQuery = params.Encode()
	}

	urlPath := Url.String()

	//请求
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//设置头
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	if len(headers) > 0 {
		for s, s2 := range headers[0] {
			req.Header.Set(s, s2)
		}
	}

	return doRequest(req)
}

func HttpPostXML(url string, data interface{}, headers ...map[string]string) ([]byte, error) {
	b, err := xml.Marshal(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	if len(headers) > 0 {
		for s, s2 := range headers[0] {
			req.Header.Set(s, s2)
		}
	}

	return doRequest(req)
}

func doRequest(req *http.Request) (body []byte, err error) {
	//超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	//客户端
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	//请求
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(`StatusCode:%d`, resp.StatusCode)
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
