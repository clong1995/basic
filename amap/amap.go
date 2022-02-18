package amap

import (
	"basic/amap/geocode"
	"basic/color"
)

var Amap *server

type (
	Server struct {
		WebKey string
	}

	server struct {
		server Server
	}
)

//Geo 地理编码 API 服务地址，名称转坐标
func (s server) Geo(address string) (res geocode.GeoRes, err error) {
	return geocode.Geo(s.server.WebKey, address)
}

//ReGeo 地理编码 API 服务地址，坐标转名称
func (s server) ReGeo(location string) (geocodes geocode.ReGeocode, err error) {
	return geocode.ReGeo(s.server.WebKey, location)
}

//ReGeoSmart 地理编码 API 服务地址，坐标转名称
func (s server) ReGeoSmart(location string) (res geocode.ReGeoSmartRes, err error) {
	return geocode.ReGeoSmart(s.server.WebKey, location)
}

//ReGeoSmartList 地理编码 API 服务地址，坐标转名称
func (s server) ReGeoSmartList(location string) (res geocode.ReGeoSmartListRes, err error) {
	return geocode.ReGeoSmartList(s.server.WebKey, location)
}

//Search 搜索
func (s server) Search(keywords, types, region string) (res []geocode.SearchPoi, err error) {
	return geocode.Search(s.server.WebKey, keywords, types, region)
}

//Detail 根据AOI或POI的id查询
func (s server) Detail(id string) (res geocode.DetailPoiRes, err error) {
	return geocode.Detail(s.server.WebKey, id)
}

func (s Server) Run() {
	//防止多次创建
	if Amap != nil {
		return
	}
	//创建对象
	Amap = &server{server: s}
	color.Success("[amap] create client success")
}
