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
func (s server) Geo(address string) (res geocode.GeoRes, err error) {
	return geocode.Geo(s.server.WebKey, address)
}

//ReGeo 地理编码 API 服务地址，名称转坐标
func (s server) ReGeo(location string) (geocodes geocode.ReGeocode, err error) {
	return geocode.ReGeo(s.server.WebKey, location)
}

//ReGeoSmart 地理编码 API 服务地址，名称转坐标
func (s server) ReGeoSmart(location string) (res geocode.ReGeoSmartRes, err error) {
	return geocode.ReGeoSmart(s.server.WebKey, location)
}

//Search
func (s server) Search(keywords, region string) (res []geocode.SearchPoi, err error) {
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
