package wechat

import (
	"basic/color"
	"basic/id"
	"basic/random"
	"basic/request"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Wx *server

var (
	ErrInvalidBlockSize    = errors.New("invalid block size")
	ErrInvalidPKCS7Data    = errors.New("invalid PKCS7 data")
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

type server struct {
	server Server
}
type wXLoginResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}
type wxUserInfo struct {
	OpenID          string `json:"openId"`
	UnionID         string `json:"unionId"`
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	NickName        string `json:"nickName"`
	Gender          int64  `json:"gender"`
	City            string `json:"city"`
	Province        string `json:"province"`
	Country         string `json:"country"`
	AvatarURL       string `json:"avatarUrl"`
	Language        string `json:"language"`
	Watermark       struct {
		Timestamp int64  `json:"timestamp"`
		AppID     string `json:"appid"`
	} `json:"watermark"`
}
type Server struct {
	AppID, SecretKey, PayKey, MchId, PayNotifyUrl string
}
type unifiedorderData struct {
	XMLName        struct{} `xml:"xml"`
	AppId          string   `xml:"appid"`
	Attach         string   `xml:"attach"`
	Body           string   `xml:"body"`
	MchId          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	NotifyUrl      string   `xml:"notify_url"`
	OpenId         string   `xml:"openid"`
	OutTradeNo     string   `xml:"out_trade_no"`
	Sign           string   `xml:"sign"`
	SpbillCreateIp string   `xml:"spbill_create_ip"`
	TotalFee       int      `xml:"total_fee"`
	TradeType      string   `xml:"trade_type"`
}

type Sign struct {
	OutTradeNo string `json:"out_trade_no"`
	TimeStamp  string `json:"time_stamp"`
	NonceStr   string `json:"nonce_str"`
	Package    string `json:"package_str"`
	PaySign    string `json:"pay_sign"`
}

// Decrypt 解密
func (s *server) Decrypt(encryptedData, sessionKey, iv string) (info *wxUserInfo, err error) {
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		log.Println(err)
		return
	}

	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		log.Println(err)
		return
	}

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		log.Println(err)
		return
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, err = pkcs7Unpad(cipherText, block.BlockSize())
	if err != nil {
		log.Println(err)
		return
	}
	info = new(wxUserInfo)
	err = json.Unmarshal(cipherText, info)
	if err != nil {
		log.Println(err)
		return
	}

	return
}
func (s *server) WXLogin(code string) (*wXLoginResp, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		s.server.AppID, s.server.SecretKey, code)

	// 创建http get请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析http请求中body 数据到我们定义的结构体中
	wxResp := wXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		return nil, err
	}

	// 判断微信接口返回的是否是一个异常情况
	if wxResp.ErrCode != 0 {
		return nil, errors.New(fmt.Sprintf("ErrCode:%s  ErrMsg:%s", wxResp.ErrCode, wxResp.ErrMsg))
	}

	return &wxResp, nil
}
func (s *server) PaySign(body, openID, spbillCreateIp string, totalFee int) (sign Sign, err error) {
	timeUnix := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := random.String(16)
	outTradeNo := id.SId.String()
	prepayId, err := unifiedorder(body, s.server.AppID, s.server.MchId, s.server.PayNotifyUrl, s.server.PayKey, openID, nonceStr, outTradeNo, spbillCreateIp, totalFee)
	if err != nil {
		log.Println(err)
		return
	}
	packageStr := "prepay_id=" + prepayId
	data := []byte(fmt.Sprintf(
		"appId=%s&nonceStr=%s&package=%s&signType=MD5&timeStamp=%s&key=%s",
		s.server.AppID,
		nonceStr,
		packageStr,
		timeUnix,
		s.server.PayKey,
	))
	has := md5.Sum(data)
	md5str := strings.ToUpper(fmt.Sprintf("%x", has))
	return Sign{
		OutTradeNo: outTradeNo,
		TimeStamp:  timeUnix,
		NonceStr:   nonceStr,
		Package:    packageStr,
		PaySign:    md5str,
	}, nil
}
func (s Server) CreateClient() {
	//防止多次创建
	if Wx != nil {
		return
	}
	//创建对象
	Wx = &server{server: s}
	color.Success(fmt.Sprintf("[wechat] create node %s success", s.AppID))
}

// pkcs7Unpad returns slice of the original data without padding
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return data[:len(data)-n], nil
}
func unifiedorder(body, appID, mchId, payNotifyUrl, payKey, openid, nonceStr, outTradeNo, spbillCreateIp string, totalFee int) (string, error) {
	attach := "BigAnt Pay"
	//step 1
	stringA := fmt.Sprintf(`appid=%s&attach=%s&body=%s&mch_id=%s&nonce_str=%s&notify_url=%s&openid=%s&out_trade_no=%s&spbill_create_ip=%s&total_fee=%d&trade_type=JSAPI`,
		appID, attach, body, mchId, nonceStr, payNotifyUrl, openid, outTradeNo, spbillCreateIp, totalFee)

	//step 2
	stringSignTemp := stringA + "&key=" + payKey
	sign := strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(stringSignTemp))))

	uData := unifiedorderData{
		AppId:          appID,
		Attach:         attach,
		Body:           body,
		MchId:          mchId,
		NonceStr:       nonceStr,
		NotifyUrl:      payNotifyUrl,
		OpenId:         openid,
		OutTradeNo:     outTradeNo,
		Sign:           sign,
		SpbillCreateIp: spbillCreateIp,
		TotalFee:       totalFee,
		TradeType:      "JSAPI",
	}
	//统一下单参数
	//https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
	//校验地址
	//https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=20_1
	data, err := request.HttpPostXML("https://api.mch.weixin.qq.com/pay/unifiedorder", uData)
	if err != nil {
		return "", err
	}

	type xmlData struct {
		ReturnCode string `xml:"return_code"`
		ReturnMsg  string `xml:"return_msg"`
		PrepayId   string `xml:"prepay_id"`
	}
	var root xmlData
	err = xml.Unmarshal(data, &root)
	if err != nil {
		return "", err
	}

	if root.ReturnCode != "SUCCESS" {
		log.Println(root)
		err = fmt.Errorf("支付失败")
		return "", err
	}

	//log.Println(root)

	return root.PrepayId, nil
}
