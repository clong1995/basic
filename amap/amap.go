package amap

import (
	"basic/amap/geocode"
	"basic/color"
)

var Amap *server

type server struct {
	server Server
}
type Server struct {
	WebKey string
}

//Geo 地理编码 API 服务地址，名称转坐标
func (s server) Geo(address string) (err error, res geocode.GeoRes) {
	return geocode.Geo(s.server.WebKey, address)
}

//ReGeo 地理编码 API 服务地址，名称转坐标
func (s server) ReGeo(location string) (err error, geocodes geocode.ReGeocode) {
	return geocode.ReGeo(s.server.WebKey, location)
}

//ReGeoSmart 地理编码 API 服务地址，名称转坐标
func (s server) ReGeoSmart(location string) (err error, res geocode.ReGeoSmartRes) {
	return geocode.ReGeoSmart(s.server.WebKey, location)
}

//Search
func (s server) Search(keywords, region string) (err error, res []geocode.SearchPoi) {
	return geocode.Search(s.server.WebKey, keywords, region)
}

func (s Server) CreateClient() {
	//防止多次创建
	if Amap != nil {
		return
	}
	//创建对象
	Amap = &server{server: s}
	color.Success("[amap] create client success")
}
