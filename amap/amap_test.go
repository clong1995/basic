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

//这个不用测试，用不到
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
	res, err := Amap.ReGeoSmart("36.308863,120.439877")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Logf("%+v", res)
}

func Test_server_ReGeoSmartList(t *testing.T) {
	Server{
		WebKey: "",
	}.Run()
	res, err := Amap.ReGeoSmartList("36.308863,120.439877")
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
	res, err := Amap.Search("拙政园", "", "苏州市")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Logf("%+v", res)
}

func Test_server_Driving(t *testing.T) {
	Server{
		WebKey: "",
	}.Run()
	_, err := Amap.Driving("31.432890,120.577894", "31.415495,120.566179")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
}
func Test_server_DrivingPolyline(t *testing.T) {
	Server{
		WebKey: "",
	}.Run()
	res, err := Amap.DrivingPolyline("39.90920999352365,116.39739791108934", "39.8819159595503,116.41078749826261")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Logf("%+v", res)
}

func Test_server_DrivingPointsPolyline(t *testing.T) {
	Server{
		WebKey: "",
	}.Run()
	res, err := Amap.DrivingPointsPolyline("39.8819159595503,116.41078749826261;39.87981240581741,116.43825331810524;39.865478942122955,116.42881194253437")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Logf("%+v", res)
}

func Test_server_ReGeoContains(t *testing.T) {
	Server{
		WebKey: "",
	}.Run()
	res, err := Amap.ReGeoContains("39.909167,116.397441", "北京市东城区")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Logf("%+v", res)
}
