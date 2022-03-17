package ws

import (
	"basic/color"
	"basic/id"
	"basic/ip"
	"basic/token"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var WS *server

type (
	OnStart   func() error
	OnClose   func(id string)
	OnMessage func(data []byte)
	OnConn    func(uid, data string) error

	Server struct {
		Addr            string    //监听地址
		ReadBufferSize  int       //最大消息长度
		WriteBufferSize int       //最大head息长度
		Origin          bool      //是否是用于web，跨域
		UserAgent       string    //允许的UserAgent
		OnStart         OnStart   //启动时
		OnClose         OnClose   //当关闭一个链接时
		OnMessage       OnMessage //当收到消息
		OnConn          OnConn    //当链接时，data的内容是链接参数的d参数
		Block           bool      //当主协程能自己维持，block不用开启
	}

	server struct {
		SendMessage func(id string, data []byte) error
	}

	client struct {
		userId  string
		conn    *websocket.Conn
		message chan []byte
	}
)

var clients map[string]*client

//var onStart OnStart
var onClose OnClose
var onMessage OnMessage

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func (s Server) Run() {
	err := s.OnStart()
	if err != nil {
		log.Println(err)
		return
	}

	//当不配置的时候，使用以下默认配置
	if s.Addr == "" {
		s.Addr = ":80"
	}
	if s.ReadBufferSize == 0 {
		s.ReadBufferSize = 1024
	}
	if s.WriteBufferSize == 0 {
		s.WriteBufferSize = 1024
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:    s.ReadBufferSize,
		WriteBufferSize:   s.WriteBufferSize,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			//return s.Origin
			return true
		},
	}

	//外部的函数
	onClose = s.OnClose
	onMessage = s.OnMessage

	clients = make(map[string]*client)

	//s:签名
	//t:token
	//sec:秒时间戳
	//d:数据
	//签名的数据为 sec+t+d
	//s=signature&t=token&sec=xxxx&d=p,c,d
	sh := func(w http.ResponseWriter, r *http.Request) {
		//解析参数
		var m = make(map[string]string)
		for key, value := range r.URL.Query() {
			m[key] = value[0]
		}
		if len(m) == 0 {
			log.Println("param is empty")
			return
		}

		//检查超过时间
		sec := ""
		if value, ok := m["sec"]; ok {
			sec = value
			var i64 int64
			i64, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				log.Println(err)
				return
			}

			if time.Now().Unix()-i64 > 15 {
				log.Println("connection time out")
				return
			}
		} else {
			log.Println("sec is empty")
			return
		}

		//签名
		signature := ""
		if value, ok := m["s"]; ok {
			signature = value
		} else {
			log.Println("s is empty")
			return
		}
		//token
		token_ := ""
		if value, ok := m["t"]; ok {
			token_ = value
		} else {
			log.Println("t is empty")
			return
		}
		//数据
		data := ""
		if value, ok := m["d"]; ok {
			data = value
		} else {
			log.Println("d is empty")
			return
		}

		//提取token信息
		tk := token.Token{}
		if err = tk.Decode(token_); err != nil {
			log.Println("decode token err:", err)
			return
		}
		userId := id.SId.ToString(tk.Id)

		//检查签名
		var buffer bytes.Buffer
		buffer.Write([]byte(sec + token_ + data))
		buffer.Write([]byte(tk.AccessKeyID()))
		sum := md5.Sum(buffer.Bytes())
		if signature != hex.EncodeToString(sum[:]) {
			log.Println("signature err")
			return
		}

		//检查是否有遗留链接断开之前的链接
		for i, c := range clients {
			if i == userId {
				c.close()
			}
		}

		//创建新的链接
		var conn *websocket.Conn
		conn, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		if err = conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Println(err)
			return
		}
		conn.SetPongHandler(func(string) error {
			//log.Println("pong handler")
			return conn.SetReadDeadline(time.Now().Add(pongWait))
		})

		c := &client{
			userId:  userId,
			conn:    conn,
			message: make(chan []byte),
		}

		//调用外部链接方法
		if err = s.OnConn(userId, data); err != nil {
			log.Println(err)
			if err = c.conn.Close(); err != nil {
				log.Println(err)
				return
			}
			return
		}

		//加入clients列表
		clients[userId] = c
		//启动每一个conn
		go c.pump()
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", sh)

	go func() {
		err = http.ListenAndServe(s.Addr, mux)
		if err != nil {
			log.Println("[ws] Listen error!", err)
			return
		}
	}()

	//ws操作句柄
	WS = new(server)
	WS.SendMessage = func(id string, data []byte) (err error) {
		if c, ok := clients[id]; ok {
			err = c.send(data)
		} else {
			err = fmt.Errorf("%s offline", id)
		}
		return
	}

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

	color.Success(fmt.Sprintf(
		"[ws] %s listening ws://%s%s ",
		s.UserAgent,
		ips[0],
		s.Addr,
	))

	if s.Block {
		select {}
	}
}

func (c *client) pump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.close()
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	//读取消息
	go func() {
		defer wg.Done()
		for {
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println(err)
				}
				break
			}
			//收到的消息
			onMessage(message)
		}
	}()

	//ping
	go func() {
		defer wg.Done()
		done := make(chan struct{}, 1)
		for {
			select {
			//ping
			case <-ticker.C:
				if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
					log.Println(err)
					close(done)
					return
				}
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Println(err)
					close(done)
					return
				}
				//log.Println("ping message")
			case <-done:
				return
			}
		}
	}()

	wg.Wait()
}

func (c *client) close() {
	//关闭
	if err := c.conn.Close(); err != nil {
		log.Println(err)
	}
	//调用关闭
	onClose(c.userId)
	//删除
	if _, ok := clients[c.userId]; ok {
		delete(clients, c.userId)
	}
	//log.Println("close conn")
}

func (c *client) send(data []byte) (err error) {
	if err = c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		log.Println(err)
		c.close()
		return
	}

	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Println(err)
		c.close()
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		c.close()
		return
	}

	if err = w.Close(); err != nil {
		log.Println(err)
		c.close()
		return
	}
	return
}
