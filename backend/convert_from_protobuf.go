// backend包提供插件后端的核心功能
package backend

// 导入所需的包
import (
	// Protobuf生成的包
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
	// 时间包
	"time"
)

// ConvertFromProtobuf 提供将Protobuf数据类型转换为SDK类型的功能
type ConvertFromProtobuf struct{}

// FromProto 创建新的从Protobuf转换器
// 返回：
//   - ConvertFromProtobuf: 从Protobuf转换器实例
func FromProto() ConvertFromProtobuf {
	return ConvertFromProtobuf{}
}

// User 将Protobuf用户对象转换为SDK用户对象
// 参数：
//   - user: Protobuf用户对象
//
// 返回：
//   - *User: SDK用户对象
func (f ConvertFromProtobuf) User(user *pluginv2.User) *User {
	if user == nil {
		return nil
	}

	return &User{
		Login: user.Login,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}

// AppInstanceSettings 将Protobuf应用实例设置转换为SDK格式
// 参数：
//   - proto: Protobuf应用实例设置
//
// 返回：
//   - *AppInstanceSettings: SDK应用实例设置
func (f ConvertFromProtobuf) AppInstanceSettings(proto *pluginv2.AppInstanceSettings) *AppInstanceSettings {
	if proto == nil {
		return nil
	}

	return &AppInstanceSettings{
		JSONData:                proto.JsonData,
		DecryptedSecureJSONData: proto.DecryptedSecureJsonData,
		Updated:                 time.Unix(0, proto.LastUpdatedMS*int64(time.Millisecond)),
	}
}

// DataSourceInstanceSettings 将Protobuf数据源实例设置转换为SDK格式
// 参数：
//   - proto: Protobuf数据源实例设置
//   - pluginID: 插件ID
//
// 返回：
//   - *DataSourceInstanceSettings: SDK数据源实例设置
func (f ConvertFromProtobuf) DataSourceInstanceSettings(proto *pluginv2.DataSourceInstanceSettings, pluginID string) *DataSourceInstanceSettings {
	if proto == nil {
		return nil
	}

	return &DataSourceInstanceSettings{
		ID:                      proto.Id,
		UID:                     proto.Uid,
		Type:                    pluginID,
		Name:                    proto.Name,
		URL:                     proto.Url,
		User:                    proto.User,
		Database:                proto.Database,
		BasicAuthEnabled:        proto.BasicAuthEnabled,
		BasicAuthUser:           proto.BasicAuthUser,
		JSONData:                proto.JsonData,
		DecryptedSecureJSONData: proto.DecryptedSecureJsonData,
		Updated:                 time.Unix(0, proto.LastUpdatedMS*int64(time.Millisecond)),
	}
}

// PluginContext 将Protobuf插件上下文转换为SDK格式
// 参数：
//   - proto: Protobuf插件上下文
//
// 返回：
//   - PluginContext: SDK插件上下文
func (f ConvertFromProtobuf) PluginContext(proto *pluginv2.PluginContext) PluginContext {
	return PluginContext{
		OrgID:                      proto.OrgId,
		PluginID:                   proto.PluginId,
		User:                       f.User(proto.User),
		AppInstanceSettings:        f.AppInstanceSettings(proto.AppInstanceSettings),
		DataSourceInstanceSettings: f.DataSourceInstanceSettings(proto.DataSourceInstanceSettings, proto.PluginId),
	}
}

// CallResourceRequest 将Protobuf资源调用请求转换为SDK格式
// 参数：
//   - protoReq: Protobuf资源调用请求
//
// 返回：
//   - *CallResourceRequest: SDK资源调用请求
func (f ConvertFromProtobuf) CallResourceRequest(protoReq *pluginv2.CallResourceRequest) *CallResourceRequest {
	headers := map[string][]string{}
	for k, values := range protoReq.Headers {
		headers[k] = values.Values
	}

	return &CallResourceRequest{
		PluginContext: f.PluginContext(protoReq.PluginContext),
		Path:          protoReq.Path,
		Method:        protoReq.Method,
		URL:           protoReq.Url,
		Headers:       headers,
		Body:          protoReq.Body,
	}
}

// CallResourceResponse 将Protobuf资源调用响应转换为SDK格式
// 参数：
//   - protoResp: Protobuf资源调用响应
//
// 返回：
//   - *CallResourceResponse: SDK资源调用响应
func (f ConvertFromProtobuf) CallResourceResponse(protoResp *pluginv2.CallResourceResponse) *CallResourceResponse {
	headers := map[string][]string{}
	for k, values := range protoResp.Headers {
		headers[k] = values.Values
	}

	return &CallResourceResponse{
		Status:  int(protoResp.Code),
		Body:    protoResp.Body,
		Headers: headers,
	}
}

// CheckHealthRequest 将Protobuf健康检查请求转换为SDK格式
// 参数：
//   - protoReq: Protobuf健康检查请求
//
// 返回：
//   - *CheckHealthRequest: SDK健康检查请求
func (f ConvertFromProtobuf) CheckHealthRequest(protoReq *pluginv2.CheckHealthRequest) *CheckHealthRequest {
	if protoReq.Headers == nil {
		protoReq.Headers = map[string]string{}
	}

	return &CheckHealthRequest{
		PluginContext: f.PluginContext(protoReq.PluginContext),
		Headers:       protoReq.Headers,
	}
}

// CheckHealthResponse 将Protobuf健康检查响应转换为SDK格式
// 参数：
//   - protoResp: Protobuf健康检查响应
//
// 返回：
//   - *CheckHealthResult: SDK健康检查结果
func (f ConvertFromProtobuf) CheckHealthResponse(protoResp *pluginv2.CheckHealthResponse) *CheckHealthResult {
	status := HealthStatusUnknown
	switch protoResp.Status {
	case pluginv2.CheckHealthResponse_ERROR:
		status = HealthStatusError
	case pluginv2.CheckHealthResponse_OK:
		status = HealthStatusOk
	}

	return &CheckHealthResult{
		Status:      status,
		Message:     protoResp.Message,
		JSONDetails: protoResp.JsonDetails,
	}
}

// Pub2ChannelRequest 将Protobuf发布频道请求转换为SDK格式
// 参数：
//   - protoReq: Protobuf发布频道请求
//
// 返回：
//   - *Pub2ChannelRequest: SDK发布频道请求
func (f ConvertFromProtobuf) Pub2ChannelRequest(protoReq *pluginv2.Pub2ChannelRequest) *Pub2ChannelRequest {
	return &Pub2ChannelRequest{
		PluginContext: f.PluginContext(protoReq.PluginContext),
		Channel:       protoReq.Channel,
		Data:          protoReq.JsonDetails,
	}
}

// Pub2ChannelResponse 将Protobuf发布频道响应转换为SDK格式
// 参数：
//   - protoResp: Protobuf发布频道响应
//
// 返回：
//   - *Pub2ChannelResponse: SDK发布频道响应
func (f ConvertFromProtobuf) Pub2ChannelResponse(protoResp *pluginv2.Pub2ChannelResponse) *Pub2ChannelResponse {
	status := PubStatusUnknown
	switch protoResp.Status {
	case pluginv2.Pub2ChannelResponse_ERROR:
		status = PubStatusError
	case pluginv2.Pub2ChannelResponse_OK:
		status = PubStatusOk
	}

	return &Pub2ChannelResponse{
		Status:      status,
		Message:     protoResp.Message,
		JsonDetails: protoResp.JsonDetails,
	}
}

// PluginInfoGetReq 将Protobuf插件信息请求转换为SDK格式
// 参数：
//   - protoReq: Protobuf插件信息请求
//
// 返回：
//   - *PluginInfoGetReq: SDK插件信息请求
func (f ConvertFromProtobuf) PluginInfoGetReq(protoReq *pluginv2.PluginInfoGetReq) *PluginInfoGetReq {

	return &PluginInfoGetReq{
		PluginContext: f.PluginContext(protoReq.PluginContext),
	}
}

// PluginInfoGetRes 将Protobuf插件信息响应转换为SDK格式
// 参数：
//   - protoResp: Protobuf插件信息响应
//
// 返回：
//   - *PluginInfoGetRes: SDK插件信息响应
func (f ConvertFromProtobuf) PluginInfoGetRes(protoResp *pluginv2.PluginInfoGetRes) *PluginInfoGetRes {

	return &PluginInfoGetRes{
		PluginID:      protoResp.PluginID,
		PluginVersion: protoResp.PluginVersion,
	}
}

// TaskRequest 将Protobuf任务请求转换为SDK格式
// 参数：
//   - protoReq: Protobuf任务请求
//
// 返回：
//   - *TaskRequest: SDK任务请求
func (f ConvertFromProtobuf) TaskRequest(protoReq *pluginv2.TaskRequest) *TaskRequest {
	return &TaskRequest{
		JSONData: protoReq.Data,
		TaskName: protoReq.TaskName,
		UserID:   int(protoReq.UserId),
	}
}

// TaskResponse 将Protobuf任务响应转换为SDK格式
// 参数：
//   - protoResp: Protobuf任务响应
//
// 返回：
//   - *TaskResponse: SDK任务响应
func (f ConvertFromProtobuf) TaskResponse(protoResp *pluginv2.TaskResponse) *TaskResponse {
	status := TaskStatusUnknown
	switch protoResp.Status {
	case pluginv2.TaskResponse_FAILED:
		status = TaskStatusFailed
	case pluginv2.TaskResponse_SUCCESS:
		status = TaskStatusSuccess
	case pluginv2.TaskResponse_RUNNING:
		status = TaskStatusRunning
	}

	return &TaskResponse{
		Status:  status,
		Message: protoResp.Message,
	}
}
