package excel

import (
	"fmt"
	"github.com/clong1995/basic/color"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"strings"
)

type (
	Server struct{}
	server struct{}
)

var (
	Excel *server
)

func (s server) AllRows(filename string) (rows [][]string, err error) {
	if filename == "" {
		err = fmt.Errorf("file null")
		log.Println(err)
		return
	}
	var f *excelize.File

	if strings.HasPrefix(filename, "http") { //网络文件
		var resp *http.Response
		resp, err = http.Get(filename)
		if err != nil {
			log.Println(err)
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		f, err = excelize.OpenReader(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		f, err = excelize.OpenFile(filename)
		if err != nil {
			log.Println(err)
			return
		}
	}

	defer func() {
		_ = f.Close()
	}()

	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		err = fmt.Errorf("sheet list empty")
		log.Println(err)
		return
	}

	rows, err = f.GetRows(sheetList[0])
	if err != nil {
		log.Println(err)
		return
	}

	if len(rows) == 0 {
		err = fmt.Errorf("empty sheet")
		log.Println(err)
		return
	}

	for _, row := range rows {
		for c, colCell := range row {
			row[c] = strings.Trim(colCell, " ")
		}
	}

	return
}

func (s Server) Run() {
	if Excel != nil {
		return
	}

	Excel = &server{}
	color.Success("[excel] create client success")
}
