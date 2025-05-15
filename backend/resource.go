// backend包提供插件后端的核心功能
package backend

// 导入所需的包
import (
	// 上下文包
	"context"
)

// CallResourceRequest 资源调用请求结构
type CallResourceRequest struct {
	// 插件上下文
	PluginContext PluginContext
	// 请求路径
	Path string
	// HTTP方法
	Method string
	// 完整URL
	URL string
	// HTTP头信息
	Headers map[string][]string
	// 请求体
	Body []byte
}

// CallResourceResponse 资源调用响应结构
type CallResourceResponse struct {
	// HTTP状态码
	Status int
	// HTTP头信息
	Headers map[string][]string
	// 响应体
	Body []byte
}

// CallResourceResponseSender 资源调用响应发送器接口
type CallResourceResponseSender interface {
	// Send 发送响应
	Send(*CallResourceResponse) error
}

// CallResourceHandler 资源调用处理器接口
type CallResourceHandler interface {
	// CallResource 处理资源调用
	// 参数：
	//   - ctx: 上下文
	//   - req: 资源调用请求
	//   - sender: 响应发送器
	//
	// 返回：
	//   - error: 错误信息
	CallResource(ctx context.Context, req *CallResourceRequest, sender CallResourceResponseSender) error
}

// CallResourceHandlerFunc 资源调用处理函数类型
type CallResourceHandlerFunc func(ctx context.Context, req *CallResourceRequest, sender CallResourceResponseSender) error

// CallResource 实现CallResourceHandler接口
func (fn CallResourceHandlerFunc) CallResource(ctx context.Context, req *CallResourceRequest, sender CallResourceResponseSender) error {
	return fn(ctx, req, sender)
}
