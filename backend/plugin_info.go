// backend包提供插件后端的核心功能
package backend

// 导入所需的包
import (
	// 上下文包
	"context"
	// HTTP包
	"net/http"
	// 数值转换包
	"strconv"
)

// CheckHealthHandler 健康检查处理器接口
type CheckHealthHandler interface {
	// CheckHealth 检查健康状态
	// 参数：
	//   - ctx: 上下文
	//   - req: 健康检查请求
	//
	// 返回：
	//   - *CheckHealthResult: 健康检查结果
	//   - error: 错误信息
	CheckHealth(ctx context.Context, req *CheckHealthRequest) (*CheckHealthResult, error)
}

// CheckHealthHandlerFunc 健康检查处理函数类型
type CheckHealthHandlerFunc func(ctx context.Context, req *CheckHealthRequest) (*CheckHealthResult, error)

// CheckHealth 实现CheckHealthHandler接口
func (fn CheckHealthHandlerFunc) CheckHealth(ctx context.Context, req *CheckHealthRequest) (*CheckHealthResult, error) {
	return fn(ctx, req)
}

// HealthStatus 健康状态枚举
type HealthStatus int

// 定义健康状态常量
const (
	// 未知状态
	HealthStatusUnknown HealthStatus = iota
	// 正常状态
	HealthStatusOk
	// 错误状态
	HealthStatusError
)

// 健康状态名称映射
var healthStatusNames = map[int]string{
	0: "UNKNOWN",
	1: "OK",
	2: "ERROR",
}

// String 将健康状态转换为字符串
// 返回：
//   - string: 状态字符串表示
func (hs HealthStatus) String() string {
	s, exists := healthStatusNames[int(hs)]
	if exists {
		return s
	}
	return strconv.Itoa(int(hs))
}

// CheckHealthRequest 健康检查请求结构
type CheckHealthRequest struct {
	// 插件上下文
	PluginContext PluginContext
	// HTTP头信息
	Headers map[string]string
}

// SetHTTPHeader 设置HTTP头
// 参数：
//   - key: 头名称
//   - value: 头值
func (req *CheckHealthRequest) SetHTTPHeader(key, value string) {
	if req.Headers == nil {
		req.Headers = map[string]string{}
	}

	setHTTPHeaderInStringMap(req.Headers, key, value)
}

// DeleteHTTPHeader 删除HTTP头
// 参数：
//   - key: 头名称
func (req *CheckHealthRequest) DeleteHTTPHeader(key string) {
	deleteHTTPHeaderInStringMap(req.Headers, key)
}

// GetHTTPHeader 获取HTTP头值
// 参数：
//   - key: 头名称
//
// 返回：
//   - string: 头值
func (req *CheckHealthRequest) GetHTTPHeader(key string) string {
	return req.GetHTTPHeaders().Get(key)
}

// GetHTTPHeaders 获取所有HTTP头
// 返回：
//   - http.Header: HTTP头对象
func (req *CheckHealthRequest) GetHTTPHeaders() http.Header {
	return getHTTPHeadersFromStringMap(req.Headers)
}

// CheckHealthResult 健康检查结果结构
type CheckHealthResult struct {
	// 健康状态
	Status HealthStatus
	// 消息内容
	Message string
	// JSON详情
	JSONDetails []byte
}

// PluginInfoHandler 插件信息处理器接口
type PluginInfoHandler interface {
	// PluginInfo 获取插件信息
	// 参数：
	//   - ctx: 上下文
	//   - req: 插件信息请求
	//
	// 返回：
	//   - *PluginInfoGetRes: 插件信息响应
	//   - error: 错误信息
	PluginInfo(ctx context.Context, req *PluginInfoGetReq) (*PluginInfoGetRes, error)
}

// PluginInfoHandlerFunc 插件信息处理函数类型
type PluginInfoHandlerFunc func(ctx context.Context, req *PluginInfoGetReq) (*PluginInfoGetRes, error)

// PluginInfo 实现PluginInfoHandler接口
func (fn PluginInfoHandlerFunc) PluginInfo(ctx context.Context, req *PluginInfoGetReq) (*PluginInfoGetRes, error) {
	return fn(ctx, req)
}

// PluginInfoGetReq 插件信息获取请求结构
type PluginInfoGetReq struct {
	// 插件上下文
	PluginContext PluginContext
}

// PluginInfoGetRes 插件信息获取响应结构
type PluginInfoGetRes struct {
	// 插件ID
	PluginID string
	// 插件版本
	PluginVersion string
}

// 确保CheckHealthRequest实现ForwardHTTPHeaders接口
var _ ForwardHTTPHeaders = (*CheckHealthRequest)(nil)
