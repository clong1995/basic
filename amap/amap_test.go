package amap

import (
	"testing"
)

func Test_server_Detail(t *testing.T) {
	Server{
		WebKey: "",
	}.Run()
	res, err := Amap.Detail("B000A837FH")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Logf("%+v", res)
}

func Test_server_Geo(t *testing.T) {
	Server{
		WebKey: "",
	}.Run()
	res, err := Amap.Geo("北京市朝阳区阜通东大街6号")
	if err != nil {
		t.Errorf(err.Error())
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
		t.Errorf(err.Error())
		return
	}
	t.Logf("%+v", res)
}

func Test_server_ReGeoSmart(t *testing.T) {
	Server{
		WebKey: "",
	}.Run()
	res, err := Amap.ReGeoSmart("31.298091,120.633647")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Logf("%+v", res)
}

func Test_server_Search(t *testing.T) {
	Server{
		WebKey: "",
	}.Run()
	res, err := Amap.Search("拙政园", "苏州市")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Logf("%+v", res)
}

//13611248094 杨
