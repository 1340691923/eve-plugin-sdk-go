// backend包提供插件后端的核心功能
package backend

// 导入所需的包
import (
	// Protobuf生成的包
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
	// 时间包
	"time"
)

// ConvertToProtobuf 提供将SDK数据类型转换为Protobuf类型的功能
type ConvertToProtobuf struct{}

// ToProto 创建新的Protobuf转换器
// 返回：
//   - ConvertToProtobuf: Protobuf转换器实例
func ToProto() ConvertToProtobuf {
	return ConvertToProtobuf{}
}

// User 将用户对象转换为Protobuf用户对象
// 参数：
//   - user: 用户对象
//
// 返回：
//   - *pluginv2.User: Protobuf用户对象
func (t ConvertToProtobuf) User(user *User) *pluginv2.User {
	if user == nil {
		return nil
	}

	return &pluginv2.User{
		Login: user.Login,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}

// AppInstanceSettings 将应用实例设置转换为Protobuf格式
// 参数：
//   - s: 应用实例设置
//
// 返回：
//   - *pluginv2.AppInstanceSettings: Protobuf应用实例设置
func (t ConvertToProtobuf) AppInstanceSettings(s *AppInstanceSettings) *pluginv2.AppInstanceSettings {
	if s == nil {
		return nil
	}

	return &pluginv2.AppInstanceSettings{
		JsonData:                s.JSONData,
		DecryptedSecureJsonData: s.DecryptedSecureJSONData,
		LastUpdatedMS:           s.Updated.UnixNano() / int64(time.Millisecond),
	}
}

// DataSourceInstanceSettings 将数据源实例设置转换为Protobuf格式
// 参数：
//   - s: 数据源实例设置
//
// 返回：
//   - *pluginv2.DataSourceInstanceSettings: Protobuf数据源实例设置
func (t ConvertToProtobuf) DataSourceInstanceSettings(s *DataSourceInstanceSettings) *pluginv2.DataSourceInstanceSettings {
	if s == nil {
		return nil
	}

	return &pluginv2.DataSourceInstanceSettings{
		Id:                      s.ID,
		Uid:                     s.UID,
		Name:                    s.Name,
		Url:                     s.URL,
		User:                    s.User,
		Database:                s.Database,
		BasicAuthEnabled:        s.BasicAuthEnabled,
		BasicAuthUser:           s.BasicAuthUser,
		JsonData:                s.JSONData,
		DecryptedSecureJsonData: s.DecryptedSecureJSONData,
		LastUpdatedMS:           s.Updated.UnixNano() / int64(time.Millisecond),
	}
}

// PluginContext 将插件上下文转换为Protobuf格式
// 参数：
//   - pluginCtx: 插件上下文
//
// 返回：
//   - *pluginv2.PluginContext: Protobuf插件上下文
func (t ConvertToProtobuf) PluginContext(pluginCtx PluginContext) *pluginv2.PluginContext {
	return &pluginv2.PluginContext{
		OrgId:                      pluginCtx.OrgID,
		PluginId:                   pluginCtx.PluginID,
		User:                       t.User(pluginCtx.User),
		AppInstanceSettings:        t.AppInstanceSettings(pluginCtx.AppInstanceSettings),
		DataSourceInstanceSettings: t.DataSourceInstanceSettings(pluginCtx.DataSourceInstanceSettings),
	}
}

// CallResourceResponse 将资源调用响应转换为Protobuf格式
// 参数：
//   - resp: 资源调用响应
//
// 返回：
//   - *pluginv2.CallResourceResponse: Protobuf资源调用响应
func (t ConvertToProtobuf) CallResourceResponse(resp *CallResourceResponse) *pluginv2.CallResourceResponse {
	headers := map[string]*pluginv2.StringList{}

	for key, values := range resp.Headers {
		headers[key] = &pluginv2.StringList{Values: values}
	}

	return &pluginv2.CallResourceResponse{
		Headers: headers,
		Code:    int32(resp.Status),
		Body:    resp.Body,
	}
}

// CallResourceRequest 将资源调用请求转换为Protobuf格式
// 参数：
//   - req: 资源调用请求
//
// 返回：
//   - *pluginv2.CallResourceRequest: Protobuf资源调用请求
func (t ConvertToProtobuf) CallResourceRequest(req *CallResourceRequest) *pluginv2.CallResourceRequest {
	protoReq := &pluginv2.CallResourceRequest{
		PluginContext: t.PluginContext(req.PluginContext),
		Path:          req.Path,
		Method:        req.Method,
		Url:           req.URL,
		Body:          req.Body,
	}
	if req.Headers == nil {
		return protoReq
	}
	protoReq.Headers = make(map[string]*pluginv2.StringList, len(protoReq.Headers))
	for k, values := range req.Headers {
		protoReq.Headers[k] = &pluginv2.StringList{Values: values}
	}
	return protoReq
}

// PluginInfoGetRes 将插件信息响应转换为Protobuf格式
// 参数：
//   - protoResp: 插件信息响应
//
// 返回：
//   - *pluginv2.PluginInfoGetRes: Protobuf插件信息响应
func (f ConvertToProtobuf) PluginInfoGetRes(protoResp *PluginInfoGetRes) *pluginv2.PluginInfoGetRes {

	return &pluginv2.PluginInfoGetRes{
		PluginID:      protoResp.PluginID,
		PluginVersion: protoResp.PluginVersion,
	}
}

// CheckHealthResponse 将健康检查结果转换为Protobuf格式
// 参数：
//   - res: 健康检查结果
//
// 返回：
//   - *pluginv2.CheckHealthResponse: Protobuf健康检查响应
func (t ConvertToProtobuf) CheckHealthResponse(res *CheckHealthResult) *pluginv2.CheckHealthResponse {
	return &pluginv2.CheckHealthResponse{
		Status:      t.HealthStatus(res.Status),
		Message:     res.Message,
		JsonDetails: res.JSONDetails,
	}
}

// HealthStatus 将健康状态转换为Protobuf健康状态
// 参数：
//   - status: 健康状态
//
// 返回：
//   - pluginv2.CheckHealthResponse_HealthStatus: Protobuf健康状态
func (t ConvertToProtobuf) HealthStatus(status HealthStatus) pluginv2.CheckHealthResponse_HealthStatus {
	switch status {
	case HealthStatusUnknown:
		return pluginv2.CheckHealthResponse_UNKNOWN
	case HealthStatusOk:
		return pluginv2.CheckHealthResponse_OK
	case HealthStatusError:
		return pluginv2.CheckHealthResponse_ERROR
	}
	panic("unsupported protobuf health status type in sdk")
}

// Pub2ChannelResponse 将发布频道响应转换为Protobuf格式
// 参数：
//   - res: 发布频道响应
//
// 返回：
//   - *pluginv2.Pub2ChannelResponse: Protobuf发布频道响应
func (t ConvertToProtobuf) Pub2ChannelResponse(res *Pub2ChannelResponse) *pluginv2.Pub2ChannelResponse {
	return &pluginv2.Pub2ChannelResponse{
		Status:      t.PubStatus(res.Status),
		Message:     res.Message,
		JsonDetails: res.JsonDetails,
	}
}

// PubStatus 将发布状态转换为Protobuf发布状态
// 参数：
//   - status: 发布状态
//
// 返回：
//   - pluginv2.Pub2ChannelResponse_PubStatus: Protobuf发布状态
func (t ConvertToProtobuf) PubStatus(status PubStatus) pluginv2.Pub2ChannelResponse_PubStatus {
	switch status {
	case PubStatusUnknown:
		return pluginv2.Pub2ChannelResponse_UNKNOWN
	case PubStatusOk:
		return pluginv2.Pub2ChannelResponse_OK
	case PubStatusError:
		return pluginv2.Pub2ChannelResponse_ERROR
	}
	panic("unsupported protobuf pub status type in sdk")
}

// TaskRequest 将任务请求转换为Protobuf格式
// 参数：
//   - req: 任务请求
//
// 返回：
//   - *pluginv2.TaskRequest: Protobuf任务请求
func (t ConvertToProtobuf) TaskRequest(req *TaskRequest) *pluginv2.TaskRequest {
	return &pluginv2.TaskRequest{
		Data:     req.JSONData,
		TaskName: req.TaskName,
		UserId:   int32(req.UserID),
	}
}

// TaskResponse 将任务响应转换为Protobuf格式
// 参数：
//   - res: 任务响应
//
// 返回：
//   - *pluginv2.TaskResponse: Protobuf任务响应
func (t ConvertToProtobuf) TaskResponse(res *TaskResponse) *pluginv2.TaskResponse {
	return &pluginv2.TaskResponse{
		Status:  t.TaskStatus(res.Status),
		Message: res.Message,
	}
}

// TaskStatus 将任务状态转换为Protobuf任务状态
// 参数：
//   - status: 任务状态
//
// 返回：
//   - pluginv2.TaskResponse_TaskStatus: Protobuf任务状态
func (t ConvertToProtobuf) TaskStatus(status TaskStatus) pluginv2.TaskResponse_TaskStatus {
	switch status {
	case TaskStatusUnknown:
		return pluginv2.TaskResponse_UNKNOWN
	case TaskStatusSuccess:
		return pluginv2.TaskResponse_SUCCESS
	case TaskStatusFailed:
		return pluginv2.TaskResponse_FAILED
	case TaskStatusRunning:
		return pluginv2.TaskResponse_RUNNING
	}
	panic("unsupported protobuf task status type in sdk")
}
