package geocode

import (
	"basic/request"
	"encoding/json"
	"errors"
	"log"
)

//逆地理编码API服务地址，坐标转名称

type reGeoSmartRes struct {
	Status    string         `json:"status"` //"1"成功
	Info      string         `json:"info"`
	ReGeocode ReGeoSmartCode `json:"regeocode"`
}

type ReGeoSmartCode struct {
	FormattedAddress string                `json:"formatted_address"` //结构化地址信息
	AddressComponent AddressComponentSmart `json:"addressComponent"`  //地址元素列表
	Roads            []RoadSmart           `json:"roads"`             //道路信息列表
	Pois             []PoiSmart            `json:"pois"`              //poi信息列表，兴趣点
	Aois             []AoiSmart            `json:"aois"`              //aoi信息列表
}

type AddressComponentSmart struct {
	Province string      `json:"province"` //省
	City     interface{} `json:"city"`     //市
	District interface{} `json:"district"` //区
	Township interface{} `json:"township"` //乡镇/街道
}

type RoadSmart struct {
	Name     string `json:"name"` //道路名称
	Location string `json:"location"`
}

type PoiSmart struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type AoiSmart struct {
	Name     string `json:"name"` //
	Location string `json:"location"`
}

type ReGeoSmartRes struct {
	FormattedAddress string    `json:"formatted_address"`
	Province         string    `json:"province"` //省
	City             string    `json:"city"`     //市
	District         string    `json:"district"` //区
	Township         string    `json:"township"` //乡镇/街道
	Road             RoadSmart `json:"road"`
	Poi              PoiSmart  `json:"poi"`
	Aoi              AoiSmart  `json:"aoi"`
}

// location "经,纬"
func ReGeoSmart(key, location string) (res ReGeoSmartRes, err error) {
	resBytes, err := request.HttpGet("https://restapi.amap.com/v3/geocode/regeo", map[string]string{
		"key":        key,
		"location":   location,
		"extensions": "all",
	})
	if err != nil {
		log.Println(err)
		return
	}
	//解析
	resp := new(reGeoSmartRes)
	err = json.Unmarshal(resBytes, resp)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.Status != "1" {
		err = errors.New(resp.Info)
		log.Println(resp)
		return
	}
	res.FormattedAddress = resp.ReGeocode.FormattedAddress
	res.Province = resp.ReGeocode.AddressComponent.Province
	if val, ok := resp.ReGeocode.AddressComponent.City.(string); ok {
		res.City = val
	}
	if val, ok := resp.ReGeocode.AddressComponent.District.(string); ok {
		res.District = val
	}
	if val, ok := resp.ReGeocode.AddressComponent.Township.(string); ok {
		res.Township = val
	}
	if len(resp.ReGeocode.Roads) > 0 {
		res.Road = resp.ReGeocode.Roads[0]
	}
	if len(resp.ReGeocode.Pois) > 0 {
		res.Poi = resp.ReGeocode.Pois[0]
	}
	if len(resp.ReGeocode.Aois) > 0 {
		res.Aoi = resp.ReGeocode.Aois[0]
	}

	return
}
