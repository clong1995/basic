package amap

import (
	"testing"
)

func Test_server_Geo(t *testing.T) {
	Server{
		WebKey: "",
	}.Run()
	res, err := Amap.Geo("北京市市辖区")
	if err != nil {
		return
	}
	t.Logf("%+v", res)
}

func Test_server_ReGeo(t *testing.T) {
	Server{
		WebKey: "",
	}.Run()
	res, err := Amap.ReGeo("116.405281,39.904987")
	if err != nil {
		return
	}
	t.Logf("%+v", res)
}

func Test_server_ReGeoSmart(t *testing.T) {
	Server{
		WebKey: "",
	}.Run()
	res, err := Amap.ReGeoSmart("120.633647,31.298091")
	if err != nil {
		return
	}
	t.Logf("%+v", res)
}

func Test_server_Search(t *testing.T) {
	Server{
		WebKey: "",
	}.Run()
	err, res := Amap.Search("拙政园", "苏州市")
	if err != nil {
		return
	}
	t.Logf("%+v", res)
}

//13611248094 杨
