package geocode

import (
	"basic/request"
	"encoding/json"
	"errors"
	"log"
)

//地理编码 API 服务地址，名称转坐标，适合精确描述，标志性，如青岛五四广场、山东省青岛市
//非标志性使用POI

type searchResp struct {
	Status string      `json:"status"` //"1"成功
	Info   string      `json:"info"`
	Pois   []SearchPoi `json:"pois"`
}
type SearchPoi struct {
	Name     string `json:"name"`
	Location string `json:"location"` //坐标
	Province string `json:"pname"`    //省
	City     string `json:"cityname"` //市
	District string `json:"adname"`   //区
	Address  string `json:"address"`  //
}

func Search(key, keywords, region string) (err error, res []SearchPoi) {
	resBytes, err := request.HttpGet("https://restapi.amap.com/v5/place/text", map[string]string{
		"key":      key,
		"keywords": keywords,
		"region":   region,
	})
	if err != nil {
		log.Println(err)
		return
	}
	//解析
	resp := new(searchResp)
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

	return nil, resp.Pois
}
