package geocode

import (
	"basic/request"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

type (
	photos struct {
		Url string `json:"url"`
	}

	//解析原始数据
	searchPoi struct {
		Name     string   `json:"name"`
		Location string   `json:"location"` //坐标
		PName    string   `json:"pname"`    //省
		CityName string   `json:"cityname"` //市
		AdName   string   `json:"adname"`   //区
		Address  string   `json:"address"`  //街道
		Photos   []photos `json:"photos"`   //图片
		Type     string   `json:"type"`     //类型
	}

	searchResp struct {
		Status string      `json:"status"` //"1"成功
		Info   string      `json:"info"`
		Pois   []searchPoi `json:"pois"`
	}

	//SearchPoi 返回数据
	SearchPoi struct {
		Province string   `json:"province"` //省
		City     string   `json:"city"`     //市
		District string   `json:"district"` //区
		Address  string   `json:"address"`  //街道
		Place    string   `json:"place"`    //地点
		Location string   `json:"location"` //坐标
		Types    []string `json:"types"`    //类型
		Photos   []string `json:"photos"`   //图片
	}
)

// Search 关键字搜索，优先采用关键字
func Search(key, keywords, types, region string) (res []SearchPoi, err error) {
	region = strings.ReplaceAll(region, "市辖区", "")
	region = strings.ReplaceAll(region, "县", "")

	param := map[string]string{
		"key":         key,
		"region":      region,
		"show_fields": "photos",
		"city_limit":  "true",
		"page_size":   "25",
	}

	if keywords == "" && types == "" {
		err = fmt.Errorf("keywords和types必选其一")
		return
	}

	if keywords != "" {
		param["keywords"] = keywords
	} else {
		param["types"] = types
	}

	resBytes, err := request.HttpGet("https://restapi.amap.com/v5/place/text", param)
	if err != nil {
		log.Println(err)
		return
	}

	//log.Println(string(resBytes))

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
			Types:    []string{},
			Photos:   []string{},
		}

		//type 去重
		if poi.Type != "" {
			typeMap := map[string]int{}
			for _, s := range strings.Split(poi.Type, ";") {
				if strings.Contains(s, "|") {
					for _, s2 := range strings.Split(s, "|") {
						typeMap[s2] = 0
					}
				} else {
					typeMap[s] = 0
				}
			}
			for k := range typeMap {
				sp.Types = append(sp.Types, k)
			}
		}

		//图片数组
		if len(poi.Photos) > 0 {
			for _, photo := range poi.Photos {
				sp.Photos = append(sp.Photos, photo.Url)
			}
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
