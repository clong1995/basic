package route

import "log"

type (
	// RoutesMap 路由存储结构
	RoutesMap map[string]Route
	// Handle 函数签名
	Handle func(string, []byte) (interface{}, error)

	// IpHandle 返回IP的签名
	IpHandle func(string, string, []byte) (interface{}, error)

	// Route 一个路由的结构
	Route struct {
		Url         string
		ContentType string
		Pattern     Pattern
		handle      Handle
		ipHandle    IpHandle
	}

	//route 路由表
	route struct {
		routes RoutesMap
	}
)

var Routes route

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
		route.Pattern.UserAgent = Enable
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
