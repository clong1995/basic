// Package nlp 自然语言处理
package nlp

import (
	"basic/color"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"log"
	"regexp"
	"strings"
)

var NLP *server

type server struct {
	accessKeyId     string
	accessKeySecret string
}

type result struct {
	RequestId string `json:"RequestId"`
	Success   bool   `json:"success"`
	Data      string `json:"data"`
}

type EntityItem struct {
	Synonym string `json:"synonym"`
	Weight  string `json:"weight"`
	Tag     string `json:"tag"`
	Word    string `json:"word"`
}

type EntityResult struct {
	Result []EntityItem `json:"result"`
}

type Sentiment struct {
	PositiveProb float64 `json:"positive_prob"`
	Sentiment    string  `json:"sentiment"`
	NegativeProb float64 `json:"negative_prob"`
}

type SentimentResult struct {
	Result Sentiment `json:"result"`
}

type Server struct {
	AccessKeyId     string
	AccessKeySecret string
}

var sentimentValue = map[string]int{
	"负面": -1,
	"中性": 0,
	"正面": 1,
}

// Entity 命名实体
func (s server) Entity(text string) (entityList []EntityItem, err error) {
	request := nlpRequest()
	request.ApiName = "GetNerChEcom"
	request.QueryParams["Text"] = text

	entityResult := new(EntityResult)
	err = s.resultData(request, entityResult)
	if err != nil {
		log.Println(err)
		return
	}
	return entityResult.Result, err
}

// Sentiment 情绪分析
func (s server) Sentiment(text string) (int, error) {
	//-1,0,1
	//负面，中性，正面
	request := nlpRequest()
	request.ApiName = "GetSaChGeneral"
	request.QueryParams["Text"] = text
	sentimentResult := new(SentimentResult)
	err := s.resultData(request, sentimentResult)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return sentimentValue[sentimentResult.Result.Sentiment], err
}

// Keywords 关键词
//例如: 三只松鼠开心果、零食坚果炒货mac M1 A19男装32a无漂白原味健康食品小吃iphone12 huawei_t800 02-1
func (s server) Keywords(text string) []string {
	entityList, err := s.Entity(text)
	if err != nil {
		log.Println(err)
		return nil
	}
	keywords := make([]string, 0)
	for _, item := range entityList {
		match := false
		//单独数字
		match, err = regexp.Match(`^[0-9]*$`, []byte(item.Word))
		if err != nil {
			log.Println(err)
			return nil
		}
		if match {
			continue
		}

		//字母下划线标点空白
		match, err = regexp.Match(`^[a-zA-Z_\s\p{P}]$`, []byte(item.Word))
		if err != nil {
			log.Println(err)
			return nil
		}
		if match {
			continue
		}

		//同义词
		if item.Synonym != "" {
			keywords = append(keywords, strings.Split(item.Synonym, ";")...)
		}
		keywords = append(keywords, item.Word)
	}
	//去重复
	set := make(map[string]struct{}, len(keywords))
	j := 0
	for _, v := range keywords {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		keywords[j] = v
		j++
	}

	return keywords[:j]
}

func (s server) resultData(request *requests.CommonRequest, v interface{}) (err error) {
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", s.accessKeyId, s.accessKeySecret)
	if err != nil {
		log.Println(err)
		return
	}
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.Println(err)
		return
	}
	if !response.IsSuccess() {
		err = fmt.Errorf(response.GetHttpContentString())
		log.Println(err)
		return
	}
	contentBytes := response.GetHttpContentBytes()
	rst := new(result)
	err = json.Unmarshal(contentBytes, rst)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal([]byte(rst.Data), v)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func nlpRequest() (req *requests.CommonRequest) {
	req = requests.NewCommonRequest()
	req.Method = "GET"
	req.Scheme = "https"
	req.Domain = "alinlp.cn-hangzhou.aliyuncs.com"
	req.Version = "2020-06-29"
	req.QueryParams["ServiceCode"] = "alinlp"
	return
}

func (s Server) Run() {
	//防止多次创建
	if NLP != nil {
		return
	}

	//创建对象
	NLP = &server{
		accessKeyId:     s.AccessKeyId,
		accessKeySecret: s.AccessKeySecret,
	}
	color.Success("[NLP] create client success")
}
