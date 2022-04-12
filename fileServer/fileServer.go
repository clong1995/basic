package fileServer

import (
	"fmt"
	"github.com/clong1995/basic/color"
	"github.com/clong1995/basic/ip"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	FileServer *server
)

type (
	Server struct {
		Block      bool //当主协程能自己维持，block不用开启
		Endpoint   string
		BucketName string
	}
	server struct {
		storagePath string
		address     string
	}
)

// PutSignFileIdURL 获取上传地址
func (s server) PutSignFileIdURL(fId string) (url string) {
	return fmt.Sprintf("%s/%s", s.address, fId)
}

func (s Server) Run() {
	//防止多次创建
	if FileServer != nil {
		return
	}
	home, dirErr := os.UserHomeDir()
	if dirErr != nil {
		log.Println(dirErr)
		return
	}
	storagePath := fmt.Sprintf("%s/%s", home, s.BucketName)

	//上传
	upload := func(w http.ResponseWriter, r *http.Request) {
		//关闭
		defer func() {
			_ = r.Body.Close()
		}()

		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		//创建目录
		filePath := fmt.Sprintf("%s%s", storagePath, r.URL.Path)
		filePathArr := strings.Split(filePath, "/")
		dirPath := strings.Join(filePathArr[:len(filePathArr)-1], "/")
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//创建文件
		//var file *os.File
		file, err := os.Create(filePath)
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//读body
		r.Body = http.MaxBytesReader(w, r.Body, 10485760) //10M
		bodyByte, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errStr := fmt.Sprintf("%s : %s\n", r.URL.Path, "读取body错误")
			log.Println(errStr)
			http.Error(w, errStr, http.StatusBadRequest)
			return
		}

		if len(bodyByte) == 0 {
			errStr := fmt.Sprintf("%s : %s\n", r.URL.Path, "没有数据")
			log.Println(errStr)
			http.Error(w, errStr, http.StatusBadRequest)
			return
		}

		//写入文件
		_, err = file.Write(bodyByte)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	//下载
	download := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		filePath := fmt.Sprintf("%s%s", storagePath, strings.TrimPrefix(r.URL.Path, "/image"))
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(fileBytes)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", upload)
	mux.HandleFunc("/image/", download)

	go func() {
		err := http.ListenAndServe(s.Endpoint, mux)
		if err != nil {
			log.Fatalln(color.Red, err, color.Reset)
			return
		}
	}()

	ips, err := ip.BoundLocalIP()
	if err != nil {
		log.Println(err)
		return
	}
	if len(ips) == 0 {
		err = fmt.Errorf("no ip")
		log.Println(err)
		return
	}

	//创建对象
	FileServer = &server{
		storagePath: storagePath,
		address:     fmt.Sprintf("http://%s%s", ips[0], s.Endpoint),
	}

	color.Success(fmt.Sprintf(
		"[fileserver] listening %s ,storage:%s",
		FileServer.address,
		storagePath,
	))

	if s.Block {
		select {}
	}
}
