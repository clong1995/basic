package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/clong1995/basic/request"
	"log"
	"strings"
)

//地理编码 API 服务地址，名称转坐标，适合精确描述，标志性，如青岛五四广场、山东省青岛市
//非标志性使用POI

type (

	//解析原始数据
	detailPoi struct {
		Name     string `json:"name"`
		Location string `json:"location"` //坐标
		PName    string `json:"pname"`    //省
		CityName string `json:"cityname"` //市
		AdName   string `json:"adname"`   //区
		Address  string `json:"address"`  //街道
	}
	detailResp struct {
		Status string      `json:"status"` //"1"成功
		Info   string      `json:"info"`
		Pois   []detailPoi `json:"pois"`
	}

	// DetailPoiRes 返回数据
	DetailPoiRes struct {
		Province string `json:"province"` //省
		City     string `json:"city"`     //市
		District string `json:"district"` //区
		Address  string `json:"address"`  //街道
		Place    string `json:"place"`    //名称
		Location string `json:"location"` //坐标
	}
)

func Detail(key, id string) (res DetailPoiRes, err error) {
	resBytes, err := request.HttpGet("https://restapi.amap.com/v3/place/detail", map[string]string{
		"key": key,
		"id":  id,
	})
	if err != nil {
		log.Println(err)
		return
	}
	//解析
	resp := new(detailResp)
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
	if len(resp.Pois) <= 0 {
		err = fmt.Errorf("pois is empty")
		log.Println(err)
		return
	}

	poi := resp.Pois[0]
	res = DetailPoiRes{
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
		res.Location = fmt.Sprintf("%s,%s", location[1], location[0])
	}

	return
}
