package regeo

import (
	"log"
	"strings"
)

// ReGeoContains 坐标是否在某一区域内。
// 如:39.909167,116.397441 是否在 北京市东城区
func ReGeoContains(key, location, address string) (result bool, err error) {
	if address == "" {
		return
	}
	smart, err := ReGeoSmart(key, location)
	if err != nil {
		log.Println(err)
		return
	}
	fullAddress := smart.Province + smart.City + smart.District + smart.Address + smart.Place
	return strings.HasPrefix(fullAddress, address), nil
}
