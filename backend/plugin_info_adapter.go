// backend包提供插件后端的核心功能
package backend

// 导入所需的包
import (
	// 上下文包
	"context"
	// Protobuf生成的包
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
)

// pluginInfoSDKAdapter 插件信息SDK适配器结构
type pluginInfoSDKAdapter struct {
	// 插件信息处理器
	pluginInfoHandler PluginInfoHandler
	// 健康检查处理器
	checkHealthHandler CheckHealthHandler
}

// newPluginInfoSDKAdapter 创建新的插件信息SDK适配器
// 参数：
//   - pluginInfoHandler: 插件信息处理器
//   - checkHealthHandler: 健康检查处理器
//
// 返回：
//   - *pluginInfoSDKAdapter: 插件信息SDK适配器实例
func newPluginInfoSDKAdapter(pluginInfoHandler PluginInfoHandler, checkHealthHandler CheckHealthHandler) *pluginInfoSDKAdapter {
	return &pluginInfoSDKAdapter{
		pluginInfoHandler:  pluginInfoHandler,
		checkHealthHandler: checkHealthHandler,
	}
}

// Get 获取插件信息
// 参数：
//   - ctx: 上下文
//   - protoReq: Protobuf插件信息请求
//
// 返回：
//   - *pluginv2.PluginInfoGetRes: Protobuf插件信息响应
//   - error: 错误信息
func (a *pluginInfoSDKAdapter) Get(ctx context.Context, protoReq *pluginv2.PluginInfoGetReq) (*pluginv2.PluginInfoGetRes, error) {
	if a.checkHealthHandler != nil {
		parsedReq := FromProto().PluginInfoGetReq(protoReq)

		res, err := a.pluginInfoHandler.PluginInfo(ctx, parsedReq)
		if err != nil {
			return nil, err
		}
		return ToProto().PluginInfoGetRes(res), nil
	}

	return &pluginv2.PluginInfoGetRes{
		PluginID:      "",
		PluginVersion: "",
	}, nil
}

// CheckHealth 检查健康状态
// 参数：
//   - ctx: 上下文
//   - protoReq: Protobuf健康检查请求
//
// 返回：
//   - *pluginv2.CheckHealthResponse: Protobuf健康检查响应
//   - error: 错误信息
func (a *pluginInfoSDKAdapter) CheckHealth(ctx context.Context, protoReq *pluginv2.CheckHealthRequest) (*pluginv2.CheckHealthResponse, error) {
	if a.checkHealthHandler != nil {
		parsedReq := FromProto().CheckHealthRequest(protoReq)

		res, err := a.checkHealthHandler.CheckHealth(ctx, parsedReq)
		if err != nil {
			return nil, err
		}
		return ToProto().CheckHealthResponse(res), nil
	}

	return &pluginv2.CheckHealthResponse{
		Status: pluginv2.CheckHealthResponse_OK,
	}, nil
}
