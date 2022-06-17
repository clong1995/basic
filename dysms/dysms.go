package dysms

import (
	"github.com/clong1995/basic/color"
	"github.com/clong1995/basic/random"

	//"basic/random"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"log"
)

type (
	Server struct {
		Dev             bool
		AccessKeyId     string
		AccessKeySecret string
	}
	server struct {
		dev bool
	}
)

var (
	Dysms       *server
	dysmsClient *dysmsapi20170525.Client
)

func (s server) Send(phone, signName, templateCode string) (code string, err error) {
	//验证码
	code = random.NumberNotZeroStart(6)

	//开发者模式
	if s.dev {
		log.Println(code)
	} else {
		//发送
		sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
			PhoneNumbers:  tea.String(phone),
			SignName:      tea.String(signName),
			TemplateCode:  tea.String(templateCode),
			TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%s\"}", code)),
		}
		var resp *dysmsapi20170525.SendSmsResponse
		resp, err = dysmsClient.SendSms(sendSmsRequest)
		if err != nil {
			println(err)
			return
		}

		if tea.StringValue(resp.Body.Code) == "OK" {

		} else {
			err = fmt.Errorf("dysms err %v", tea.StringValue(resp.Body.Message))
			return
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

	Dysms = &server{dev: s.Dev}
	color.Success("[dysms] create client success")
}
