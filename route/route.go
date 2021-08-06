package route

import "log"

var Routes route

// RoutesMap 路由存储结构
type RoutesMap map[string]Route

//路由表
type route struct {
	routes RoutesMap
}

// Put 向路由表注册路由
func (r route) Put(route Route) {
	if route.handle == nil && route.ipHandle == nil {
		//存在，结束程序
		log.Panicf("'%s' handle is nill", route.Url)
	}
	//检查是否存在路由
	if _, ok := r.routes[route.Url]; ok {
		//存在，结束程序
		log.Panicf("'%s' redeclared in this gateway", route.Url)
	}
	if route.Pattern.Auth == None { // 认证
		// 默认token认证
		route.Pattern.Auth = Enable
	}
	if route.Pattern.Encrypt == None { // 高级加密
		// 默认不使用高级加密
		route.Pattern.Encrypt = EncryptDisable
	}
	if route.Pattern.UserAgent == None { // User-Agent
		// 默认不实用User-Agent
		route.Pattern.UserAgent = UserAgentDisable
	}
	if route.Pattern.General == None { // 通用模式
		//默认不使用通用模式
		route.Pattern.General = GeneralDisable
	}
	r.routes[route.Url] = route
}

// All 返回路由表
func (r route) All() RoutesMap {
	return r.routes
}

// Handle 函数签名
type Handle func(string, []byte) (interface{}, error)

// IpHandle 返回IP的签名
type IpHandle func(string, string, []byte) (interface{}, error)

// Route 一个路由的结构
type Route struct {
	Url      string
	Pattern  Pattern
	handle   Handle
	ipHandle IpHandle
}

func (r Route) Register(handle Handle) {
	r.handle = handle
	Routes.Put(r)
}

func (r Route) IpRegister(ipHandle IpHandle) {
	r.ipHandle = ipHandle
	Routes.Put(r)
}

func (r Route) Handle() Handle {
	return r.handle
}

func (r Route) IpHandle() IpHandle {
	return r.ipHandle
}

func init() {
	Routes = route{RoutesMap{}}
}
