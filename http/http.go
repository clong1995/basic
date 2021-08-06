package http

import (
	"basic/cipher"
	"basic/color"
	"basic/id"
	"basic/ip"
	. "basic/route"
	"basic/token"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/time/rate"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//原始数据
/*type original struct {
	Data      string `json:"d"`
	Signature string `json:"s"`
}*/

var contentHmac = "Content-Hmac"

type auth struct {
	Token    string `json:"t"`
	DeviceId string `json:"d"`
}

//返回数据
type response struct {
	State string      `json:"state"`
	Data  interface{} `json:"data"`
}

type Server struct {
	Addr            string        //监听地址
	MaxPayloadBytes int           //最大消息长度
	MaxHeaderBytes  int           //最head息长度
	Every           time.Duration //速度(毫秒)
	Bursts          int           //流量(个)
	ReadTimeout     time.Duration //读超时
	WriteTimeout    time.Duration //写超时
	Web             bool          //是否是用于web
	UserAgent       string        //允许的UserAgent
}

// Start 启动服务
func (h Server) Start() {

	mux := http.NewServeMux()

	//限流器
	limiter := rate.NewLimiter(rate.Every(h.Every), h.Bursts)

	//执行路由表
	for s, r := range Routes.All() {
		//闭包保存路由
		func(pattern string, route Route) {
			mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
				//关闭
				defer func() {
					_ = r.Body.Close()
				}()

				//限流
				err := limiter.Wait(context.Background())
				if err != nil {
					log.Printf("%s : %s\n", pattern, err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if h.Web == true {
					//跨域
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					w.Header().Set("Access-Control-Allow-Methods", "*")
					w.Header().Set("Access-Control-Allow-Headers", "*")
					//去掉缓存
					w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
					w.Header().Set("Pragma", "no-cache")
					w.Header().Set("Expires", "0")

					//跨域侦测
					if r.Method == http.MethodOptions {
						w.WriteHeader(http.StatusOK)
						return
					}
				}

				//处理header
				header := r.Header
				if h.UserAgent != "" && route.Pattern.UserAgent == Enable {
					acc := header.Get("User-Agent")
					if acc != h.UserAgent {
						errStr := fmt.Sprintf("%s : %s\n", pattern, "User-Agent 错误")
						http.Error(w, errStr, http.StatusForbidden)
						return
					}
				}

				//检查签名
				sig := header.Get(contentHmac)
				if route.Pattern.Auth == Enable {
					//有认证必须要校验签名
					if sig == "" {
						errStr := fmt.Sprintf("%s : %s\n", pattern, "缺少数据签名")
						http.Error(w, errStr, http.StatusForbidden)
						return
					}
				}

				//paramByte期待的数据结构是
				//head:{
				//	"Content-Hmac":"signature",
				//}
				//body: {
				//	"t":"token",
				//	"d":"deviceId",
				//	"aaa":"bbb",
				//	"ccc":"ddd"
				//}
				var paramByte []byte
				//根据方法不同处理参数
				if r.Method == http.MethodGet { //TODO get没有测试
					var m = make(map[string]string)
					for key, value := range r.URL.Query() {
						m[key] = value[0]
					}
					//TODO 论证这里会不会有错
					paramByte, err = json.Marshal(m)
					if err != nil {
						errStr := fmt.Sprintf("%s : %s\n", pattern, "读取url参数错误")
						http.Error(w, errStr, http.StatusInternalServerError)
						return
					}
				} else {
					//读body
					r.Body = http.MaxBytesReader(w, r.Body, int64(h.MaxPayloadBytes))
					paramByte, err = ioutil.ReadAll(r.Body)
					if err != nil {
						errStr := fmt.Sprintf("%s : %s\n", pattern, "读取body错误")
						http.Error(w, errStr, http.StatusRequestEntityTooLarge)
						return
					}
				}

				var tId int64
				var ak []byte
				if route.Pattern.Auth == Enable { //启用认证
					if len(paramByte) > 0 {
						//提取 token、deviceId
						a := &auth{}
						err = json.Unmarshal(paramByte, a)
						if err != nil {
							errStr := fmt.Sprintf("%s : %s\n", pattern, err)
							http.Error(w, errStr, http.StatusInternalServerError)
							return
						}

						if a.Token == "" {
							errStr := fmt.Sprintf("%s : %s\n", pattern, "缺少令牌")
							http.Error(w, errStr, http.StatusNotAcceptable)
							return
						}

						//提起令牌内容
						tk := token.Token{}
						err = tk.Decode(a.Token)
						if err != nil {
							errStr := fmt.Sprintf("%s : %s\n", pattern, "令牌错误")
							http.Error(w, errStr, http.StatusNotAcceptable)
							return
						}

						tId = tk.Id
						ak = []byte(tk.AccessKeyID())

						//校验签名
						if !cipher.CheckHmacSha256(paramByte, sig, ak) {
							errStr := fmt.Sprintf("%s : %s\n", pattern, "指纹检验失败")
							http.Error(w, errStr, http.StatusNotAcceptable)
							return
						}
					} else {
						errStr := fmt.Sprintf("%s : %s\n", pattern, "body为空")
						http.Error(w, errStr, http.StatusNoContent)
						return
					}
				}

				//var jsonErr error

				//执行
				var result interface{}
				//检查是否有特殊的handle
				ipHandle := route.IpHandle()
				if ipHandle != nil {
					result, err = ipHandle(ip.XRealIp(r), id.SId.ToString(tId), paramByte)
				} else {
					result, err = route.Handle()(id.SId.ToString(tId), paramByte)
				}

				// 通用不格式直接写出
				if route.Pattern.General == Enable {
					if err != nil {
						errStr := fmt.Sprintf("%s : %s\n", pattern, err)
						http.Error(w, errStr, http.StatusInternalServerError)
						return
					}

					if result == nil {
						w.WriteHeader(http.StatusOK)
						return
					}

					//判断是bytes
					switch value := result.(type) {
					case []byte:
					default:
						errStr := fmt.Sprintf("%v is not []byte or []uint8", value)
						http.Error(w, errStr, http.StatusInternalServerError)
						return
					}

					w.WriteHeader(http.StatusOK)
					_, err = w.Write(result.([]byte))

					/*var buf bytes.Buffer
					enc := gob.NewEncoder(&buf)
					err = enc.Encode(result)
					if err != nil {
						log.Printf("%s : %s\n", pattern, err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					_, err = w.Write(buf.Bytes())*/

					if err != nil {
						log.Printf("%s : %s\n", pattern, err)
						return
					}
					return
				}

				// 处理结果
				var jsonBytes []byte
				if err != nil {
					fmt.Printf("%s : %s", pattern, err)
					jsonBytes, err = json.Marshal(response{
						err.Error(),
						nil,
					})
				} else {
					jsonBytes, err = json.Marshal(response{
						"OK",
						result,
					})
				}

				//json错误
				if err != nil {
					errStr := fmt.Sprintf("%s : %s\n", pattern, err)
					http.Error(w, errStr, http.StatusInternalServerError)
					return
				}

				//TODO 判断是否使用gzip

				//计算hmac
				if route.Pattern.Auth == Enable {
					responseSig := cipher.HmacSha256(jsonBytes, ak)
					//写入header
					w.Header().Set(contentHmac, responseSig)
				}

				//JSON
				w.WriteHeader(http.StatusOK)

				//写出结果
				_, err = w.Write(jsonBytes)
				if err != nil {
					log.Printf("%s : %s\n", pattern, err)
					return
				}

				//写入缓存
				/*if route.Pattern.Cache != None {

				}*/
			})
		}(s, r)
	}

	//当不配置的时候，使用以下默认配置
	if h.Addr == "" {
		h.Addr = ":80"
	}
	if h.MaxPayloadBytes == 0 {
		h.MaxPayloadBytes = 1 << 20
	}
	if h.MaxHeaderBytes == 0 {
		h.MaxHeaderBytes = 1 << 20
	}
	if h.Every == 0 {
		h.Every = 10 * time.Millisecond
	}
	if h.Bursts == 0 {
		h.Every = 2
	}
	if h.ReadTimeout == 0 {
		h.ReadTimeout = 10 * time.Second
	}
	if h.WriteTimeout == 0 {
		h.ReadTimeout = 10 * time.Second
	}
	color.Success(fmt.Sprintf("[http] %s listening %s", h.UserAgent, h.Addr))
	//启动服务
	server := &http.Server{
		Addr:           h.Addr,
		ReadTimeout:    h.ReadTimeout,
		WriteTimeout:   h.WriteTimeout,
		MaxHeaderBytes: h.MaxHeaderBytes,
		Handler:        mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Println("[http] Listen error!")
	}
}
