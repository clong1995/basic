package verify

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

func Unmarshal(bytes []byte, model interface{}) error {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		err := fmt.Errorf("must pass a pointer, not a value, to FieldScan destination")
		log.Println(err)
		return err
	}
	if len(bytes) == 0 {
		return review(model, "")
	}
	//转化json为struct
	if err := json.Unmarshal(bytes, model); err != nil {
		log.Println(err)
		return err
	}
	return review(model, "")
}

// 支持
// struct{}，适用于普通但参数的简单请求，没有嵌套数据
// struct{[]struct{}}，适用于批量插入的请求，有嵌套多层但每层结构固定的数据
// unmarshal:格式化后的json对象
// field:用于递归判断嵌套的json对象
// TODO 以后在完善其他类型
func review(unmarshal interface{}, supTagJson string, field ...*reflect.Value) error {
	var err error

	var items reflect.Value
	if len(field) == 1 {
		items = *field[0]
	} else {
		items = reflect.ValueOf(unmarshal).Elem()
	}

	kind := items.Type().Kind()

	switch kind {
	//入口
	case reflect.Struct:
		for i := 0; i < items.NumField(); i++ {
			valueField := items.Field(i)
			typeField := items.Type().Field(i)
			tagRequired := typeField.Tag.Get("required")

			//log.Println(name, "|", typ, "|", tagJson, "|", tagRequired)
			//log.Println("================")

			tagJson := typeField.Tag.Get("json")
			name := typeField.Name
			typ := typeField.Type.Name()

			if tagRequired == "true" {
				if typeField.Type.Name() != "" { //基本类型，存在required，string不得为空，number不得为0
					if valueField.IsZero() || valueField.Interface() == "" {
						err = fmt.Errorf("%s %s [%s:%s] 不得为空\n", supTagJson, tagJson, name, typ)
						break
					}
				} else { //非基本类型，Slice、Map等
					//log.Println("递归")
					if valueField.Type().Kind() == reflect.Slice || valueField.Type().Kind() == reflect.Map { //判断长度
						if valueField.Len() == 0 {
							err = fmt.Errorf("%s %s [%s:%s] 不得为空\n", supTagJson, tagJson, name, typ)
							break
						}
					}

					if err = review(nil, tagJson, &valueField); err != nil {
						break
					}
				}
			}
		}
	case reflect.Slice:
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			value := reflect.Indirect(item)
			if value.Kind() == reflect.Struct || value.Kind() == reflect.Slice || value.Kind() == reflect.Map {
				if err = review(nil, supTagJson, &value); err != nil {
					break
				}
			}
		}
	case reflect.Map:
		keys := items.MapKeys()
		for _, k := range keys {
			v := items.MapIndex(k)
			value := reflect.Indirect(v)
			if value.Kind() == reflect.Struct || value.Kind() == reflect.Slice || value.Kind() == reflect.Map {
				if err = review(nil, supTagJson, &value); err != nil {
					break
				}
			}
		}
	default:
		err = fmt.Errorf(`case default:  暂不支持%s`, kind)
		break
	}

	//处理struct
	//var items = reflect.ValueOf(unmarshal).Elem()

	return err
}
