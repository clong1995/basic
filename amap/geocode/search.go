package geocode

import (
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
	//解析原始数据
	searchPoi struct {
		Name     string `json:"name"`
		Location string `json:"location"` //坐标
		PName    string `json:"pname"`    //省
		CityName string `json:"cityname"` //市
		AdName   string `json:"adname"`   //区
		Address  string `json:"address"`  //街道
	}
	searchResp struct {
		Status string      `json:"status"` //"1"成功
		Info   string      `json:"info"`
		Pois   []searchPoi `json:"pois"`
	}

	//SearchPoi 返回数据
	SearchPoi struct {
		Province string `json:"province"` //省
		City     string `json:"city"`     //市
		District string `json:"district"` //区
		Address  string `json:"address"`  //街道
		Place    string `json:"place"`
		Location string `json:"location"` //坐标
	}
)

func Search(key, keywords, region string) (res []SearchPoi, err error) {
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
	for _, poi := range resp.Pois {
		sp := SearchPoi{
			Province: poi.PName,
			City:     poi.CityName,
			District: poi.AdName,
			Address:  poi.Address,
			Place:    poi.Name,
			Location: "",
		}
		//修正"经,纬"为"纬,经"
		if poi.Location != "" {
			location := strings.Split(poi.Location, ",")
			if len(location) != 2 {
				err = fmt.Errorf("location error")
				return
			}
			sp.Location = fmt.Sprintf("%s,%s", location[1], location[0])
		}
		res = append(res, sp)
	}

	return
}
