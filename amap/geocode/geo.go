package geocode

import (
	"basic/fieldCopy"
	"basic/request"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

//地理编码 API 服务地址，名称转坐标，适合精确描述，标志性，如青岛五四广场、山东省青岛市
//非标志性使用POI

type (
	GeoRes struct {
		FormattedAddress string `json:"formatted_address"`
		Province         string `json:"province"` //省
		City             string `json:"city"`     //市
		District         string `json:"district"` //区
		Township         string `json:"township"` //乡镇/街道
		Street           string `json:"street"`   //道路
		Number           string `json:"number"`   //门牌
		Location         string `json:"location"` //坐标
	}
	geoResp struct {
		Status   string     `json:"status"` //"1"成功
		Info     string     `json:"info"`
		Geocodes []geocodes `json:"geocodes"`
	}
	geocodes struct {
		FormattedAddress string      `json:"formatted_address"`
		Province         string      `json:"province"`                   //省
		City             string      `json:"city"`                       //市
		District         interface{} `json:"district" deepcopier:"skip"` //区
		Township         interface{} `json:"township" deepcopier:"skip"` //乡镇/街道
		Street           interface{} `json:"street" deepcopier:"skip"`   //道路
		Number           interface{} `json:"number" deepcopier:"skip"`   //门牌
		Location         string      `json:"location"`                   //坐标 返回值为"经,纬",这是错误的形式，应该为"纬,经"
	}
)

func Geo(key, address string) (res GeoRes, err error) {
	resBytes, err := request.HttpGet("https://restapi.amap.com/v3/geocode/geo", map[string]string{
		"key":     key,
		"address": address,
	})
	if err != nil {
		log.Println(err)
		return
	}
	//解析
	resp := new(geoResp)
	err = json.Unmarshal(resBytes, resp)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.Status != "1" {
		err = errors.New(resp.Info)
		log.Println(err)
		return
	}

	if len(resp.Geocodes) > 0 {
		r := resp.Geocodes[0]
		err = fieldCopy.FieldFrom(&res, r)
		if err != nil {
			log.Println(err)
			return
		}
		//修正District类型
		if val, ok := r.District.(string); ok {
			res.District = val
		}
		//修正Township类型
		if val, ok := r.Township.(string); ok {
			res.Township = val
		}
		//修正Street类型
		if val, ok := r.Street.(string); ok {
			res.Street = val
		}
		//修正Number类型
		if val, ok := r.Number.(string); ok {
			res.Number = val
		}
		//修正"经,纬"为"纬,经"
		if r.Location != "" {
			location := strings.Split(r.Location, ",")
			res.Location = fmt.Sprintf("%s,%s", location[1], location[0])
		}
		return
	}
	err = errors.New("结果为空")
	return
}
