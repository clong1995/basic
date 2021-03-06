// Package route 路由模式配置
package route

type (
	PatternType int64
	Pattern     struct {
		Auth        PatternType //认证
		Cache       PatternType //缓存
		CacheExpire int64       //缓存保留时间单位秒，当Cache开启的时候有效
		Encrypt     PatternType //加密
		UserAgent   PatternType //user-agent
		General     PatternType //通用模式
		Version     int64       //内部版本
	}
)

const (
	None PatternType = iota // 开始生成枚举值, 默认为0
	// Enable 开启
	Enable
	// AuthDisable 认证
	AuthDisable
	// CacheDisable 缓存
	CacheDisable
	// EncryptDisable 加密
	EncryptDisable
	// UserAgentDisable User-Agent
	UserAgentDisable
	// GeneralDisable 通用模式
	GeneralDisable
)

func (p PatternType) String() string {
	switch p {
	case None:
		return "None"

	//认证
	case CacheDisable:
		return "关闭缓存"

	//认证
	case AuthDisable:
		return "关闭认证"

	//加密
	case EncryptDisable:
		return "关闭加密加密"

	//User-Agent
	case UserAgentDisable:
		return "关闭User-Agent"
	//	通用模式
	case GeneralDisable:
		return "关闭通用模式"
	}

	return "N/A"
}
