// backend包提供插件后端的核心功能
package backend

// 导入所需的包
import (
	// 上下文包
	"context"
	// 数值转换包
	"strconv"
)

// PubStatus 发布状态枚举
type PubStatus int

// 定义发布状态常量
const (
	// 未知状态
	PubStatusUnknown PubStatus = iota
	// 成功状态
	PubStatusOk
	// 错误状态
	PubStatusError
)

// 发布状态名称映射
var pubStatusNames = map[int]string{
	0: "UNKNOWN",
	1: "OK",
	2: "ERROR",
}

// Pub2ChannelRequest 发布到频道请求结构
type Pub2ChannelRequest struct {
	// 插件上下文
	PluginContext PluginContext
	// 频道名称
	Channel string
	// 数据内容
	Data []byte
}

// Pub2ChannelResponse 发布到频道响应结构
type Pub2ChannelResponse struct {
	// 发布状态
	Status PubStatus
	// 消息内容
	Message string
	// JSON详情
	JsonDetails []byte
}

// String 将发布状态转换为字符串
// 返回：
//   - string: 状态字符串表示
func (hs PubStatus) String() string {
	s, exists := pubStatusNames[int(hs)]
	if exists {
		return s
	}
	return strconv.Itoa(int(hs))
}

// LiveHandler 实时处理器接口
type LiveHandler interface {
	// Pub2Channel 发布消息到频道
	// 参数：
	//   - ctx: 上下文
	//   - req: 发布请求
	//
	// 返回：
	//   - *Pub2ChannelResponse: 发布响应
	//   - error: 错误信息
	Pub2Channel(ctx context.Context, req *Pub2ChannelRequest) (*Pub2ChannelResponse, error)
}

// LiveHandlerFunc 实时处理函数类型
type LiveHandlerFunc func(ctx context.Context, req *Pub2ChannelRequest) (*Pub2ChannelResponse, error)

// Pub2Channel 实现LiveHandler接口
func (fn LiveHandlerFunc) Pub2Channel(ctx context.Context, req *Pub2ChannelRequest) (*Pub2ChannelResponse, error) {
	return fn(ctx, req)
}
