package geocode

import (
	"basic/request"
	"encoding/json"
	"errors"
	"log"
)

//逆地理编码API服务地址，坐标转名称

type (
	ReGeocode struct {
		FormattedAddress string           `json:"formatted_address"` //结构化地址信息
		AddressComponent AddressComponent `json:"addressComponent"`  //地址元素列表
		Roads            []Road           `json:"roads"`             //道路信息列表
		Roadinters       []Roadinter      `json:"roadinters"`        //道路交叉口列表
		Pois             []Poi            `json:"pois"`              //poi信息列表，兴趣点
		Aois             []Aoi            `json:"aois"`              //aoi信息列表
	}

	AddressComponent struct {
		Province      string          `json:"province"`      //省
		City          interface{}     `json:"city"`          //市
		District      interface{}     `json:"district"`      //区
		Township      interface{}     `json:"township"`      //乡镇/街道
		Neighborhood  Neighborhood    `json:"neighborhood"`  //社区信息列表
		Building      Building        `json:"building"`      //楼信息列表
		StreetNumber  StreetNumber    `json:"streetNumber"`  //门牌信息列表
		SeaArea       string          `json:"seaArea"`       //所属海域信息
		BusinessAreas []BusinessAreas `json:"businessAreas"` //经纬度所属商圈列表
	}

	Neighborhood struct {
		Name interface{} `json:"name"`
		Type interface{} `json:"type"`
	}

	Building struct {
		Name interface{} `json:"name"`
		Type interface{} `json:"type"`
	}

	StreetNumber struct {
		Street    string `json:"street"` //街道名称
		Number    string `json:"number"` //门牌号
		Location  string `json:"location"`
		Direction string `json:"direction"` //坐标点所处街道方位
		Distance  string `json:"distance"`  //门牌地址到请求坐标的距离
	}

	BusinessAreas struct {
		BusinessArea string `json:"businessArea"` //商圈信息
		Name         string `json:"name"`         //商圈中心点经纬度
		Location     string `json:"location"`     //商圈中心点经纬度
		Id           string `json:"id"`           //商圈所在区域的adcode
	}

	Road struct {
		Id        string `json:"id"`        //道路id
		Name      string `json:"name"`      //道路名称
		Distance  string `json:"distance"`  //道路到请求坐标的距离
		Direction string `json:"direction"` //输入点和此路的相对方位
		Location  string `json:"location"`
	}

	Roadinter struct {
		Distance   string `json:"distance"`  //道路到请求坐标的距离
		Direction  string `json:"direction"` //输入点和此路的相对方位
		Location   string `json:"location"`
		FirstId    string `json:"first_id"` //第一条道路id
		FirstName  string `json:"first_name"`
		SecondId   string `json:"second_id"` //第一条道路id
		SecondName string `json:"second_name"`
	}

	Poi struct {
		Id           string      `json:"id"`        //
		Name         string      `json:"name"`      //
		Type         string      `json:"type"`      //
		Tel          interface{} `json:"tel"`       //
		Address      interface{} `json:"address"`   //poi地址信息
		Distance     string      `json:"distance"`  //道路到请求坐标的距离
		Direction    string      `json:"direction"` //输入点和此路的相对方位
		Location     string      `json:"location"`
		Businessarea string      `json:"businessarea"` //poi所在商圈名称
	}
	Aoi struct {
		Id       string `json:"id"`   //
		Name     string `json:"name"` //
		Location string `json:"location"`
		Area     string `json:"area"`     //所属aoi点面积
		Distance string `json:"distance"` //道路到请求坐标的距离
	}

	reGeoRes struct {
		Status    string    `json:"status"` //"1"成功
		Info      string    `json:"info"`
		ReGeocode ReGeocode `json:"regeocode"`
	}
)

func ReGeo(key, location string) (reGeocodes ReGeocode, err error) {
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
	resp := new(reGeoRes)
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
	return resp.ReGeocode, nil
}
