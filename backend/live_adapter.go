// backend包提供插件后端的核心功能
package backend

// 导入所需的包
import (
	// 上下文包
	"context"
	// Protobuf生成的包
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
)

// liveSDKAdapter 实时SDK适配器结构
type liveSDKAdapter struct {
	// 实时处理器
	liveHandler LiveHandler
}

// newLiveSDKAdapter 创建新的实时SDK适配器
// 参数：
//   - handler: 实时处理器
//
// 返回：
//   - *liveSDKAdapter: 实时SDK适配器实例
func newLiveSDKAdapter(handler LiveHandler) *liveSDKAdapter {
	return &liveSDKAdapter{
		liveHandler: handler,
	}
}

// liveResponseSenderFunc 实时响应发送函数类型
type liveResponseSenderFunc func(resp *Pub2ChannelResponse) error

// Send 实现实时响应发送接口
// 参数：
//   - resp: 发布频道响应
//
// 返回：
//   - error: 错误信息
func (fn liveResponseSenderFunc) Send(resp *Pub2ChannelResponse) error {
	return fn(resp)
}

// Pub2Channel 实现Protobuf发布频道接口
// 参数：
//   - ctx: 上下文
//   - protoReq: Protobuf发布请求
//
// 返回：
//   - *pluginv2.Pub2ChannelResponse: Protobuf发布响应
//   - error: 错误信息
func (a *liveSDKAdapter) Pub2Channel(ctx context.Context, protoReq *pluginv2.Pub2ChannelRequest) (*pluginv2.Pub2ChannelResponse, error) {
	if a.liveHandler != nil {
		parsedReq := FromProto().Pub2ChannelRequest(protoReq)

		res, err := a.liveHandler.Pub2Channel(ctx, parsedReq)
		if err != nil {
			return nil, err
		}
		return ToProto().Pub2ChannelResponse(res), nil
	}

	return &pluginv2.Pub2ChannelResponse{
		Status: pluginv2.Pub2ChannelResponse_OK,
	}, nil
}
