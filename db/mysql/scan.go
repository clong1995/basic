package mysql

import (
	"basic/cipher"
	"basic/fieldCopy"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"time"
)

// FieldScan 扫描一行数据
func FieldScan(rows *sql.Rows, field interface{}) (err error) {
	//取一行
	cols, err := rows.Columns()
	if err != nil {
		log.Println(err)
		return
	}

	//判断是结构体
	valueOfModule := reflect.ValueOf(field)
	if valueOfModule.Kind() != reflect.Ptr {
		err = fmt.Errorf("must pass a pointer, not a value, to FieldScan destination")
		return
	}

	//扫描结果
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for index := range columns {
		columnPointers[index] = &columns[index]
	}
	if err = rows.Scan(columnPointers...); err != nil {
		log.Println(err)
		return
	}

	//形成key-value
	mapValue := make(map[string]interface{}, len(cols))
	for i, colName := range cols {
		mapValue[colName] = *columnPointers[i].(*interface{})
	}

	valueOfModule = reflect.Indirect(valueOfModule)

	typeOfModule := valueOfModule.Type()

	for i := 0; i < valueOfModule.NumField(); i++ {
		//定义的tag
		fieldName := typeOfModule.Field(i).Tag.Get("field")

		//TODO 嵌入的复合结构体没法处理

		//包含要处理的field
		if item, ok := mapValue[fieldName]; ok {
			vfi := valueOfModule.Field(i)
			if vfi.CanSet() && item != nil {
				//检查模型和数据的类型是否对应
				modelType := vfi.Kind().String()

				//module的属性名
				name := typeOfModule.Field(i).Name

				//检查是否是键，键的话转成字符串，主键格式化
				if typeOfModule.Field(i).Tag.Get("key") == "true" {
					//是键，检查模型是否是string
					if modelType != "string" {
						err = fmt.Errorf(`model %s is key mast be string,current type is %s`, name, modelType)
						log.Println(err)
						return
					}
					//检查是int
					switch item.(type) {
					case int64:
						vfi.SetString(cipher.Base64EncryptInt64(item.(int64)))
						break
					default:
						return fmt.Errorf(`table %s is key mast be BIGINT`, fieldName)
					}
					continue
				}
				//检查是否是时间格式化
				if typeOfModule.Field(i).Tag.Get("dateFormat") == "true" {
					//是键，检查模型是否是string
					if modelType != "string" {
						return fmt.Errorf(`model %s is dateFormat mast be string,current type is %s`, name, modelType)
					}
					//检查是time
					switch item.(type) {
					case time.Time:
						fmtData := item.(time.Time).Format("2006-01-02 15:04:05.000")
						if fmtData == "1000-01-01 00:00:00.000" {
							fmtData = ""
						}
						vfi.SetString(fmtData)
						break
					default:
						return fmt.Errorf(`table %s is dateFormat mast be datatime`, fieldName)
					}
					continue
				}

				//item 从 map 取出的 interface{}
				//log.Println(modelType)
				switch val := item.(type) {
				case int:
					if modelType != "int" {
						return modelErr(name, modelType, "int", fieldName, val)
					}
					vfi.SetInt(int64(item.(int)))
					break
				case uint:
					if modelType != "uint" {
						return modelErr(name, modelType, "uint", fieldName, val)
					}
					vfi.SetUint(uint64(item.(uint)))
					break
				case int8:
					if modelType != "int8" {
						return modelErr(name, modelType, "int8", fieldName, val)
					}
					vfi.SetInt(int64(item.(int8)))
					break
				case uint8:
					if modelType != "uint8" {
						return modelErr(name, modelType, "uint8", fieldName, val)
					}
					vfi.SetUint(uint64(item.(uint8)))
					break
				case int16:
					if modelType != "int16" {
						return modelErr(name, modelType, "int16", fieldName, val)
					}
					vfi.SetInt(int64(item.(int16)))
					break
				case uint16:
					if modelType != "uint16" {
						return modelErr(name, modelType, "uint16", fieldName, val)
					}
					vfi.SetUint(uint64(item.(uint16)))
					break
				case int32:
					if modelType != "int32" {
						return modelErr(name, modelType, "int32", fieldName, val)
					}
					vfi.SetInt(int64(item.(int32)))
					break
				case uint32:
					if modelType != "uint32" {
						return modelErr(name, modelType, "uint32", fieldName, val)
					}
					vfi.SetUint(uint64(item.(uint32)))
					break
				case int64:
					if modelType == "int64" {
						vfi.SetInt(item.(int64))
					} else if modelType == "bool" {
						vfi.Set(reflect.ValueOf(item.(int64) != 0))
					} else {
						return modelErr(name, modelType, "int64/bool", fieldName, val)
					}
					break
				case uint64:
					if modelType != "uint64" {
						return modelErr(name, modelType, "uint64", fieldName, val)
					}
					vfi.SetUint(uint64(item.(int64)))
					break
				case float32:
					if modelType != "float32" {
						return modelErr(name, modelType, "float32", fieldName, val)
					}
					vfi.SetFloat(float64(item.(float32)))
					break
				case float64:
					if modelType != "float64" {
						return modelErr(name, modelType, "float64", fieldName, val)
					}
					vfi.SetFloat(item.(float64))
					break
				case string:
					if modelType != "string" {
						return modelErr(name, modelType, "string", fieldName, val)
					}
					vfi.SetString(item.(string))
					break
				case bool:
					if modelType != "bool" {
						return modelErr(name, modelType, "bool", fieldName, val)
					}
					itemBool := item.(bool)
					vfi.Set(reflect.ValueOf(&itemBool))
					break
				case []uint8:
					if modelType != "string" {
						return modelErr(name, modelType, "string", fieldName, val)
					}
					vfi.SetString(b2s(item.([]uint8)))
					break
				case complex64:
					if modelType != "complex64" {
						return modelErr(name, modelType, "complex64", fieldName, val)
					}
					itemComplex64 := item.(complex64)
					vfi.Set(reflect.ValueOf(&itemComplex64))
					break
				case complex128:
					if modelType != "complex128" {
						return modelErr(name, modelType, "complex128", fieldName, val)
					}
					itemComplex128 := item.(complex128)
					vfi.Set(reflect.ValueOf(&itemComplex128))
					break
				case time.Time: //原始时间格式
					if vfi.Type().String() != "time.Time" {
						return modelErr(name, modelType, "time.Time", fieldName, val)
					}
					itemTime := item.(time.Time)
					vfi.Set(reflect.ValueOf(itemTime))
					break
				default:
					return modelErr(name, modelType, "unKnow", fieldName, val)
				}
			} else {
				log.Println("vfi can`t set or item is nil")
			}
		}
	}

	return
}

// FieldScanList 扫描多行数据，这里需要优化扫描列表
func FieldScanList() {

}

// RespScan 拷贝到 response
func RespScan(rows *sql.Rows, field, resp interface{}) (err error) {
	err = FieldScan(rows, field)
	if err != nil {
		log.Println(err)
		return
	}

	err = fieldCopy.FieldFrom(resp, field)
	if err != nil {
		log.Println(err)
		return err
	}
	return
}

func b2s(bs []uint8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = v
	}
	return string(b)
}

func modelErr(name, modelType, dataType, field string, val interface{}) error {
	return fmt.Errorf("model %s type is %s, table %s type is %s. value is %v", name, modelType, field, dataType, val)
}
