package jverify

import (
	"basic/cipher"
	"basic/color"
	"basic/request"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

var Jverify *server

type server struct {
	server Server
}
type Server struct {
	AppKey       string
	MasterSecret string
	PrivateKey   []byte
}
type req struct {
	LoginToken string `json:"loginToken"` //认证SDK获取到的loginToken
	ExID       string `json:"exID"`       //开发者自定义的id，非必填
}
type res struct {
	Id      int64  `json:"id"`
	Code    int64  `json:"code"`
	Content string `json:"content"`
	ExID    string `json:"exID"`
	Phone   string `json:"phone"`
}

//Decrypt 解密loginToken获取手机号
func (s *server) Decrypt(loginToken string) (phone string, err error) {
	if loginToken == "" {
		err = errors.New("loginToken 为空")
		return
	}
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.server.AppKey, s.server.MasterSecret)))

	resBytes, err := request.HttpPostJson("https://api.verification.jpush.cn/v1/web/loginTokenVerify", req{
		LoginToken: loginToken,
	}, map[string]string{
		"Authorization": "Basic " + auth,
	})
	if err != nil {
		log.Println(err)
		return
	}
	resData := new(res)
	err = json.Unmarshal(resBytes, resData)
	if err != nil {
		log.Println(err)
		return
	}
	if resData.Code != 8000 {
		err = errors.New(resData.Content)
		log.Println(err)
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(resData.Phone)
	if err != nil {
		log.Println(err)
		return
	}
	decrypt, err := cipher.RSADecrypt(decoded, s.server.PrivateKey)
	if err != nil {
		log.Println(err)
		return
	}
	return string(decrypt), nil
}

func (s Server) Run() {
	//防止多次创建
	if Jverify != nil {
		return
	}
	//创建对象
	Jverify = &server{server: s}
	color.Success("[jverify] create client success")
}
