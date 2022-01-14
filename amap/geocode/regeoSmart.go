package geocode

import (
	"basic/request"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

//逆地理编码API服务地址，坐标转名称

type (

	//解析原始数据
	reGeoSmartRes struct {
		Status    string         `json:"status"` //"1"成功
		Info      string         `json:"info"`
		ReGeocode reGeoSmartCode `json:"regeocode"`
	}

	reGeoSmartCode struct {
		FormattedAddress string                `json:"formatted_address"` //结构化地址信息
		AddressComponent addressComponentSmart `json:"addressComponent"`  //地址元素列表
		Roads            []roadSmart           `json:"roads"`             //道路信息列表
		Pois             []poiSmart            `json:"pois"`              //poi信息列表，兴趣点
	}

	addressComponentSmart struct {
		Province string      `json:"province"` //省
		City     interface{} `json:"city"`     //市
		District interface{} `json:"district"` //区
		Township interface{} `json:"township"` //乡镇/街道
	}

	roadSmart struct {
		Name     string `json:"name"` //道路名称
		Location string `json:"location"`
	}

	poiSmart struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		Location string `json:"location"`
	}

	// ReGeoSmartRes 返回数据
	ReGeoSmartRes struct {
		FormattedAddress string `json:"formatted_address"`
		Province         string `json:"province"` //省
		City             string `json:"city"`     //市
		District         string `json:"district"` //区
		Township         string `json:"township"` //乡镇/街道
		Address          string `json:"address"`
		Place            string `json:"place"`
		Location         string `json:"location"`
	}
)

// ReGeoSmart location "纬,经"
func ReGeoSmart(key, location string) (res ReGeoSmartRes, err error) {

	arr := strings.Split(location, ",")
	if len(arr) == 2 {
		location = fmt.Sprintf("%s,%s", arr[1], arr[0])
	} else {
		err = fmt.Errorf("location error")
		return
	}

	//location "经,纬"
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
	//POI
	if len(resp.ReGeocode.Pois) > 0 {
		poi := resp.ReGeocode.Pois[0]
		res.Address = poi.Address
		res.Place = poi.Name
		arr = strings.Split(poi.Location, ",")
		if len(arr) == 2 {
			res.Location = fmt.Sprintf("%s,%s", arr[1], arr[0])
		}
	} else {
		//Road
		if len(resp.ReGeocode.Roads) > 0 {
			road := resp.ReGeocode.Roads[0]
			res.Place = road.Name
			arr = strings.Split(road.Location, ",")
			if len(arr) == 2 {
				res.Location = fmt.Sprintf("%s,%s", arr[1], arr[0])
			}
		}
	}

	return
}
