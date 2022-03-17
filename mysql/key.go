package mysql

import (
	"basic/id"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Key 获取一个string key 或者 用int转
func Key(key ...int64) string {
	if len(key) == 1 {
		//转化key
		return id.SId.ToString(key[0])
	}
	return id.SId.String()
}

// DBKey 获取一个int key或者用 string转
func DBKey(key ...string) int64 {
	if len(key) == 1 {
		//转化key
		if len(key[0]) < 10 {
			log.Println(fmt.Sprintf("id格式错误: %s", key[0]))
			return 0
		}
		return id.SId.ToInt(key[0])
	}
	//生成key
	return id.SId.Int()
}

// FIND_IN_SET
//将 ["ABCATkUZFRM","ABBAWnf3FRM","ABDAp9_3FRM","ABCA3BUCFhM"] 的id列表转换为
//"1375033046692007936,1375277353218871296,1375277801195704320,1375289029125214208"
//用于 FIND_IN_SET 语句
func FIND_IN_SET(stringKeyList []string) string {
	var ids []string

	for _, sk := range stringKeyList {
		ids = append(ids, strconv.FormatInt(DBKey(sk), 10))
	}

	return strings.Join(ids, ",")
}

// REGEXP_IN_SET 用于 REGEXP 语句
func REGEXP_IN_SET(stringKeyList []string) string {
	return strings.Join(stringKeyList, "|")
}
