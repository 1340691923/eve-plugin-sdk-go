// backend包提供插件后端的核心功能
package backend

// 导入所需的包
import (
	// Protobuf生成的包
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
	// HTTP包
	"net/http"
)

// resourceSDKAdapter 资源SDK适配器结构
type resourceSDKAdapter struct {
	// 资源调用处理器
	callResourceHandler CallResourceHandler
}

// newResourceSDKAdapter 创建新的资源SDK适配器
// 参数：
//   - handler: 资源调用处理器
//
// 返回：
//   - *resourceSDKAdapter: 资源SDK适配器实例
func newResourceSDKAdapter(handler CallResourceHandler) *resourceSDKAdapter {
	return &resourceSDKAdapter{
		callResourceHandler: handler,
	}
}

// callResourceResponseSenderFunc 资源调用响应发送函数类型
type callResourceResponseSenderFunc func(resp *CallResourceResponse) error

// Send 实现CallResourceResponseSender接口
// 参数：
//   - resp: 资源调用响应
//
// 返回：
//   - error: 错误信息
func (fn callResourceResponseSenderFunc) Send(resp *CallResourceResponse) error {
	return fn(resp)
}

// CallResource 实现Protobuf资源调用服务接口
// 参数：
//   - protoReq: Protobuf资源调用请求
//   - protoSrv: Protobuf资源调用服务器
//
// 返回：
//   - error: 错误信息
func (a *resourceSDKAdapter) CallResource(protoReq *pluginv2.CallResourceRequest, protoSrv pluginv2.Resource_CallResourceServer) error {
	if a.callResourceHandler == nil {
		return protoSrv.Send(&pluginv2.CallResourceResponse{
			Code: http.StatusNotImplemented,
		})
	}

	fn := callResourceResponseSenderFunc(func(resp *CallResourceResponse) error {
		return protoSrv.Send(ToProto().CallResourceResponse(resp))
	})

	return a.callResourceHandler.CallResource(protoSrv.Context(), FromProto().CallResourceRequest(protoReq), fn)
}
