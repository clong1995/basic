package dysms

import (
	"basic/color"
	"basic/random"

	//"basic/random"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"time"
)

type codeMsg struct {
	code       string
	expiration time.Time
}

type server struct {
}

type Server struct {
	AccessKeyId     string
	AccessKeySecret string
}

var Dysms *server
var dysmsClient *dysmsapi20170525.Client
var dict = make(map[string]codeMsg)

func (s server) Send(phone, signName, templateCode string, showCode bool) (err error) {
	now := time.Now()

	//删除过期的
	for p, msg := range dict {
		if now.After(msg.expiration) {
			delete(dict, p)
		}
	}

	//构造新的
	mm, err := time.ParseDuration("5m")
	if err != nil {
		return
	}

	//验证码
	code := random.NumberNotZeroStart(6)

	if showCode {
		log.Println(code)
	}
	//发送
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%s\"}", code)),
	}
	resp, err := dysmsClient.SendSms(sendSmsRequest)
	if err != nil {
		println(err)
		return
	}

	if tea.StringValue(resp.Body.Code) == "OK" {

	} else {
		err = fmt.Errorf("dysms err %v", resp.Body.Message)
		return
	}

	//保存
	dict[phone] = codeMsg{
		code:       code,
		expiration: now.Add(mm),
	}

	return
}

func (s server) Check(phone, code string) (result bool) {
	now := time.Now()
	//删除过期的
	for p, msg := range dict {
		if now.After(msg.expiration) {
			delete(dict, p)
		}
	}

	if value, ok := dict[phone]; ok {
		//判断验证码
		if value.code == code {
			//正确
			delete(dict, phone)
			return true
		}
	}
	return
}

func (s Server) Run() {
	if Dysms != nil {
		return
	}
	var err error

	cof := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: &s.AccessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: &s.AccessKeySecret,
	}

	// 访问的域名
	cof.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	dysmsClient = &dysmsapi20170525.Client{}
	dysmsClient, err = dysmsapi20170525.NewClient(cof)
	if err != nil {
		log.Fatalln(color.Red, err, color.Reset)
	}

	Dysms = new(server)
	color.Success("[dysms] create client success")
}
