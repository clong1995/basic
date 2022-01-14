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

	// DetailPoi 返回数据
	DetailPoi struct {
		Province string `json:"province"` //省
		City     string `json:"city"`     //市
		District string `json:"district"` //区
		Address  string `json:"address"`  //街道
		Place    string `json:"place"`    //名称
		Location string `json:"location"` //坐标
	}
)

func Detail(key, id string) (res []DetailPoi, err error) {
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

	//规范输出
	//res = []DetailPoi{}
	for _, poi := range resp.Pois {
		dp := DetailPoi{
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
			dp.Location = fmt.Sprintf("%s,%s", location[1], location[0])
		}
		res = append(res, dp)
	}

	return
}
