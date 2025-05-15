// live包提供实时数据处理功能
package live

// 导入所需的包
import (
	// 导入上下文包
	"context"
	// 导入错误处理包
	"errors"
	// 导入backend核心功能包
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	// 导入JSON处理库
	"github.com/goccy/go-json"
)

// LiveHandle 定义实时处理接口
type LiveHandle interface {
	// Pub2Channel 发布消息到指定频道
	Pub2Channel(ctx context.Context, channel string, req []byte) (res map[string]interface{}, err error)
}

// Live 实现实时数据处理功能
type Live struct {
	// 实时处理器
	liveHandle LiveHandle
}

// NewLive 创建一个新的实时数据处理实例
func NewLive(liveHandle LiveHandle) *Live {
	// 返回初始化的Live结构体
	return &Live{liveHandle: liveHandle}
}

// Pub2Channel 实现向频道发布消息的功能
func (this *Live) Pub2Channel(ctx context.Context, req *backend.Pub2ChannelRequest) (*backend.Pub2ChannelResponse, error) {
	// 检查实时处理器是否存在
	if this.liveHandle == nil {
		// 如果不存在，返回错误状态和错误信息
		return &backend.Pub2ChannelResponse{Status: backend.PubStatusError}, errors.New("该插件没有实现长连接处理器")
	}

	// 调用实时处理器发布消息
	res, err := this.liveHandle.Pub2Channel(ctx, req.Channel, req.Data)
	// 检查是否有错误
	if err != nil {
		// 如果有错误，返回错误状态和错误信息
		return &backend.Pub2ChannelResponse{Status: backend.PubStatusError}, err
	}
	// 将结果序列化为JSON
	resB, err := json.Marshal(res)
	// 检查序列化是否有错误
	if err != nil {
		// 如果有错误，返回错误状态和错误信息
		return &backend.Pub2ChannelResponse{Status: backend.PubStatusError}, err
	}
	// 返回成功状态和JSON详情
	return &backend.Pub2ChannelResponse{Status: backend.PubStatusOk, JsonDetails: resB}, nil
}
