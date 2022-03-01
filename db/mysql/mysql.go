package mysql

import (
	"basic/color"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"strings"
)

type (
	Server struct {
		DataSource string
		MaxOpen    int
	}
	server struct {
	}
)

var (
	Mysql   *server
	mysqlDB *sql.DB
)

// TxBegin 开启事物
func (s server) TxBegin() (*sql.Tx, error) {
	return mysqlDB.Begin()
}

// TxEnd 关闭事物
func (s server) TxEnd(tx *sql.Tx, err error) {
	if tx == nil {
		return
	}
	var endErr error
	if err != nil {
		endErr = tx.Rollback()
	} else {
		endErr = tx.Commit()
	}
	if endErr != nil {
		log.Println(endErr)
	}
}

// RowsCloseAndTxEnd 关闭row并关闭事物
func (s server) RowsCloseAndTxEnd(rows *sql.Rows, tx *sql.Tx, err error) {
	if rows != nil {
		if cErr := rows.Close(); cErr != nil {
			log.Println(cErr)
		}
	}
	s.TxEnd(tx, err)
}

//TxExecProc 执行一条sql
func (s server) TxExecProc(tx *sql.Tx, procName string, args ...interface{}) (sql.Result, error) {
	ret, values := argsData(args)
	sqlQuery := fmt.Sprintf("CALL %s (%s)", procName, ret)
	return tx.Exec(sqlQuery, values...)
}

// TxExecMultiProc TODO 批量执行sql
/*func (s server) TxExecMultiProc(tx *sql.Tx, procName string, arr [][]interface{}) (sql.Result, error) {
	return nil, nil
}*/

// TxQueryProc 查询
func (s server) TxQueryProc(tx *sql.Tx, procName string, args ...interface{}) (*sql.Rows, error) {
	ret, values := argsData(args)
	sqlQuery := fmt.Sprintf("CALL %s (%s)", procName, ret)
	return tx.Query(sqlQuery, values...)
}

// Run 连接数据库
func (s Server) Run() {
	//防止多次创建
	if Mysql != nil {
		return
	}

	var sqlErr error
	mysqlDB, sqlErr = sql.Open("mysql", s.DataSource+"?charset=utf8mb4&loc=Asia%2FShanghai&parseTime=true&multiStatements=true")
	if sqlErr != nil {
		log.Fatalln(color.Red, sqlErr, color.Reset)
	}

	mysqlDB.SetMaxOpenConns(s.MaxOpen)
	sqlErr = mysqlDB.Ping()

	if sqlErr != nil {
		log.Fatal(color.Red, sqlErr, color.Reset)
	}
	Mysql = new(server)
	color.Success(fmt.Sprintf("[mysql] connect %s success", strings.Split(s.DataSource, "@tcp")[1]))
}

//格式化参数
func argsData(args []interface{}) (sqlArgs string, values []interface{}) {
	for _, arg := range args {
		//TODO 切片
		//结构体
		t := reflect.TypeOf(arg)
		if t.Kind() == reflect.Struct {
			v := reflect.ValueOf(arg)
			for i := 0; i < v.NumField(); i++ {
				name := t.Field(i).Name
				/*if strings.HasSuffix(name, "Id") {//排除结尾为Id的字段
					continue
				}*/
				if name == "Id" { //排除Id的字段
					continue
				}
				values = append(values, v.Field(i).Interface())
			}
		} else {
			values = append(values, arg)
		}
	}
	sqlArgs = strings.Repeat("?,", len(values))
	sqlArgs = strings.TrimRight(sqlArgs, ",")
	return
}
