package route

import "log"

type (
	// RouteMap 路由存储结构
	routeMap map[string]Route
	// Handle 函数签名。id,数据
	Handle func(string, []byte) (interface{}, error)

	// IpHandle 返回IP的签名。ip,id,数据
	IpHandle func(string, string, []byte) (interface{}, error)

	// SessionHandle 返回session的签名。session,数据
	SessionHandle func(string, []byte) (interface{}, error)

	// Route 一个路由的结构
	Route struct {
		Url           string
		ContentType   string
		Pattern       Pattern
		handle        Handle
		ipHandle      IpHandle
		sessionHandle SessionHandle
	}
)

//var Routes route
var routes routeMap

// Put 向路由表注册路由
func (r routeMap) put(route Route) {
	if route.handle == nil && route.ipHandle == nil {
		//存在，结束程序
		log.Panicf("'%s' handle is nill", route.Url)
	}
	//检查是否存在路由
	if _, ok := r[route.Url]; ok {
		//存在，结束程序
		log.Panicf("'%s' redeclared in this gateway", route.Url)
	}
	if route.Pattern.Auth == None { // 认证
		// 默认token认证
		route.Pattern.Auth = Enable
	}
	if route.Pattern.Cache == None { // 缓存
		// 默认不缓存
		route.Pattern.Cache = CacheDisable
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
		// 默认不使用通用模式
		route.Pattern.General = GeneralDisable
	}
	r[route.Url] = route
}

// All 返回路由表
func All() map[string]Route {
	return routes
}

func (r Route) Register(handle Handle) {
	r.handle = handle
	routes.put(r)
}

func (r Route) IpRegister(ipHandle IpHandle) {
	r.ipHandle = ipHandle
	routes.put(r)
}

func (r Route) SessionRegister(sessionHandle SessionHandle) {
	r.sessionHandle = sessionHandle
	routes.put(r)
}

func (r Route) Handle() Handle {
	return r.handle
}

func (r Route) IpHandle() IpHandle {
	return r.ipHandle
}

func (r Route) SessionHandle() SessionHandle {
	return r.sessionHandle
}

func init() {
	routes = routeMap{}
}
