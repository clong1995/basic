// Package route 路由模式配置
package route

type (
	Pattern struct {
		Auth      PatternType //认证
		Encrypt   PatternType //加密
		UserAgent PatternType //user-agent
		General   PatternType //通用模式
	}
	PatternType int
)

const (
	None PatternType = iota // 开始生成枚举值, 默认为0
	// Enable 开启
	Enable
	// AuthDisable 认证
	AuthDisable
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
