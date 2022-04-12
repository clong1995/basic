package direction

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/clong1995/basic/request"
	"log"
	"strings"
)

type (
	DrivingStep struct {
		Instruction  string      `json:"instruction"`   //行驶指示
		Orientation  string      `json:"orientation"`   //进入道路方向
		StepDistance string      `json:"step_distance"` //分段距离信息
		RoadName     interface{} `json:"road_name"`     //分段道路名称
		Polyline     string      `json:"polyline"`      //分路段坐标点串
	}

	DrivingPath struct {
		Distance string        `json:"distance"` //距离米
		Steps    []DrivingStep `json:"steps"`
	}

	DrivingRoute struct {
		Paths []DrivingPath `json:"paths"`
	}

	DrivingResp struct {
		Status string       `json:"status"` //"1"成功
		Info   string       `json:"info"`
		Route  DrivingRoute `json:"route"`
	}

	DrivingPolylines [][2]string //[[经度,纬度],[经度,纬度],[经度,纬度]...]
)

// Driving 行车基础信息
func Driving(key, origin, destination string) (res DrivingResp, err error) {
	if origin == "" || destination == "" {
		return
	}
	//起点
	originArr := strings.Split(origin, ",")
	origin = ""
	if len(originArr) == 2 {
		origin = fmt.Sprintf("%s,%s", originArr[1], originArr[0])
	}
	//终点
	destArr := strings.Split(destination, ",")
	destination = ""
	if len(destArr) == 2 {
		destination = fmt.Sprintf("%s,%s", destArr[1], destArr[0])
	}

	resBytes, err := request.HttpGet("https://restapi.amap.com/v5/direction/driving", map[string]string{
		"key":         key,
		"origin":      origin,
		"destination": destination,
		"show_fields": "polyline",
	})
	if err != nil {
		log.Println(err)
		return
	}

	//err = utils.WriteLog(resBytes, "/Users/yuchenglong/Desktop")

	if err != nil {
		log.Println(err)
		return
	}
	//解析
	//resp := new(DrivingResp)
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		log.Println(err)
		return
	}

	if res.Status != "1" {
		err = errors.New(res.Info)
		log.Println(err)
		return
	}

	return
}

// DrivingPolyline 返回起点和终点的Polyline
func DrivingPolyline(key, origin, destination string) (res DrivingPolylines, err error) {
	driving, err := Driving(key, origin, destination)
	if err != nil {
		log.Println(err)
		return
	}
	paths := driving.Route.Paths
	if len(paths) == 0 {
		return
	}

	steps := paths[0].Steps
	if len(steps) == 0 {
		return
	}

	//起点
	originArr := strings.Split(origin, ",")
	res = append(res, [2]string{
		originArr[1],
		originArr[0],
	})
	//中间点
	for _, step := range steps {
		polyline := step.Polyline
		if polyline == "" {
			continue
		}
		polylineArr := strings.Split(polyline, ";")
		for _, pol := range polylineArr {
			if pol != "" {
				polArr := strings.Split(pol, ",")
				if len(polArr) == 2 {
					res = append(res, [2]string{polArr[0], polArr[1]})
				}
			}
		}
	}

	//终点
	destinationArr := strings.Split(destination, ",")
	res = append(res, [2]string{
		destinationArr[1],
		destinationArr[0],
	})
	return
}

// DrivingPointsPolyline 返回一串坐标的Polyline,[points]参数用,和;隔开
func DrivingPointsPolyline(key, points string) (res DrivingPolylines, err error) {
	pointList := strings.Split(points, ";")
	if len(pointList) < 2 {
		//少于两个点
		return
	}

	for i := range pointList {
		start := i
		end := i + 1
		if end == len(pointList) { //到了结尾
			break
		}
		var polyline DrivingPolylines
		polyline, err = DrivingPolyline(key, pointList[start], pointList[end])
		if err != nil {
			log.Println(err)
			return
		}
		res = append(res, polyline...)
	}
	return
}
