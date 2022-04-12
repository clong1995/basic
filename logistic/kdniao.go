package logistic

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/clong1995/basic/color"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type (
	Trace struct {
		AcceptStation string `json:"AcceptStation"`
		AcceptTime    string `json:"AcceptTime"`
	}
	Server struct {
		EBusinessID  string
		ApiKey       string
		CallbackAddr string
	}
	tracesResult struct {
		LogisticCode string  `json:"LogisticCode"`
		ShipperCode  string  `json:"ShipperCode"`
		Traces       []Trace `json:"Traces"`
		State        string  `json:"State"`
		EBusinessID  string  `json:"EBusinessID"`
		Success      bool    `json:"Success"`
	}
	subscribeResult struct {
		ShipperCode string `json:"ShipperCode"`
		UpdateTime  string `json:"UpdateTime"`
		EBusinessID string `json:"EBusinessID"`
		Success     bool   `json:"Success"`
	}
	server struct {
		eBusinessID string
		apiKey      string
	}
)

const (
	traceUrl     string = "https://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"
	subscribeUrl string = "https://api.kdniao.com/api/dist"
)

var Logistic *server

// Traces 实时查询接口
func (s server) Traces(shipperCode, LogisticCode string) ([]Trace, error) {
	// 组装应用级参数
	RequestData := "{" +
		"'CustomerName': ''," +
		"'OrderCode': ''," +
		"'ShipperCode': '" + shipperCode + "'," +
		"'LogisticCode': '" + LogisticCode + "'," +
		"}"

	dataSign, err := getSign(RequestData, s.apiKey)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// 组装系统级参数
	v := map[string]string{
		"RequestType": "1002",
		"EBusinessID": s.eBusinessID,
		"DataType":    "2",
		"RequestData": RequestData,
		"DataSign":    dataSign,
	}
	bytes, err := post(traceUrl, v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//log.Println(string(bytes))
	res := new(tracesResult)
	err = json.Unmarshal(bytes, res)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	ts := res.Traces
	//逆向结果
	for i, j := 0, len(ts)-1; i < j; i, j = i+1, j-1 {
		ts[i], ts[j] = ts[j], ts[i]
	}
	return ts, nil
}

// Subscribe 订阅接口
func (s server) Subscribe(shipperCode, LogisticCode, Callback string) error {
	// 组装应用级参数
	RequestData := "{" +
		"'Callback':'" + Callback + "'," +
		"'ShipperCode': '" + shipperCode + "'," +
		"'LogisticCode': '" + LogisticCode + "'," +
		//"'CustomerName':'1234'," + //CustomerName字段: ShipperCode为SF时必填，对应寄件人/收件人手机号后四位；ShipperCode为其他快递时，可不填或保留字段，不可传值
		"}"
	dataSign, err := getSign(RequestData, s.apiKey)
	if err != nil {
		log.Println(err)
		return err
	}
	// 组装系统级参数
	v := map[string]string{
		"RequestType": "1008",
		"EBusinessID": s.eBusinessID,
		"DataType":    "2",
		"RequestData": RequestData,
		"DataSign":    dataSign,
	}
	bytes, err := post(subscribeUrl, v)
	if err != nil {
		log.Println(err)
		return err
	}
	res := new(subscribeResult)
	err = json.Unmarshal(bytes, res)
	if err != nil {
		log.Println(err)
		return err
	}
	if res.Success != true {
		err = fmt.Errorf(string(bytes))
		log.Println(err)
		return err
	}
	return nil
}

func base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func getSign(n, apiKey string) (string, error) {
	str := n + apiKey
	w := md5.New()
	_, err := io.WriteString(w, str)
	if err != nil {
		log.Println(err)
		return "", err
	}
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	debyte := base64Encode([]byte(md5str))
	return fmt.Sprintf("%s", debyte), nil
} //签名

func post(url string, params map[string]string) ([]byte, error) {
	var values []string
	for k, v := range params {
		values = append(values, fmt.Sprintf("%s=%s", k, v))
	}
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(strings.Join(values, "&")))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("StatusCode == %d", resp.StatusCode)
		log.Println(err)
		return nil, err
	}
	contentBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return contentBytes, nil
}

func (s Server) Run() {
	//防止多次创建
	if Logistic != nil {
		return
	}
	//创建对象
	Logistic = &server{
		eBusinessID: s.EBusinessID,
		apiKey:      s.ApiKey,
	}
	color.Success(fmt.Sprintf("[logistic] kdniao create client success,callback addr %s,\n<快递鸟>免费版仅支持 申通、圆通、百世、天天,有效期半年.", s.CallbackAddr))
}
