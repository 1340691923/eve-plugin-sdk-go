package backend

import (
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
	"time"
)

type ConvertFromProtobuf struct{}

func FromProto() ConvertFromProtobuf {
	return ConvertFromProtobuf{}
}

func (f ConvertFromProtobuf) User(user *pluginv2.User) *User {
	if user == nil {
		return nil
	}

	return &User{
		Login:	user.Login,
		Name:	user.Name,
		Email:	user.Email,
		Role:	user.Role,
	}
}

func (f ConvertFromProtobuf) AppInstanceSettings(proto *pluginv2.AppInstanceSettings) *AppInstanceSettings {
	if proto == nil {
		return nil
	}

	return &AppInstanceSettings{
		JSONData:	proto.JsonData,
		DecryptedSecureJSONData:	proto.DecryptedSecureJsonData,
		Updated:	time.Unix(0, proto.LastUpdatedMS*int64(time.Millisecond)),
	}
}

func (f ConvertFromProtobuf) DataSourceInstanceSettings(proto *pluginv2.DataSourceInstanceSettings, pluginID string) *DataSourceInstanceSettings {
	if proto == nil {
		return nil
	}

	return &DataSourceInstanceSettings{
		ID:	proto.Id,
		UID:	proto.Uid,
		Type:	pluginID,
		Name:	proto.Name,
		URL:	proto.Url,
		User:	proto.User,
		Database:	proto.Database,
		BasicAuthEnabled:	proto.BasicAuthEnabled,
		BasicAuthUser:	proto.BasicAuthUser,
		JSONData:	proto.JsonData,
		DecryptedSecureJSONData:	proto.DecryptedSecureJsonData,
		Updated:	time.Unix(0, proto.LastUpdatedMS*int64(time.Millisecond)),
	}
}

func (f ConvertFromProtobuf) PluginContext(proto *pluginv2.PluginContext) PluginContext {
	return PluginContext{
		OrgID:	proto.OrgId,
		PluginID:	proto.PluginId,
		User:	f.User(proto.User),
		AppInstanceSettings:	f.AppInstanceSettings(proto.AppInstanceSettings),
		DataSourceInstanceSettings:	f.DataSourceInstanceSettings(proto.DataSourceInstanceSettings, proto.PluginId),
	}
}

func (f ConvertFromProtobuf) CallResourceRequest(protoReq *pluginv2.CallResourceRequest) *CallResourceRequest {
	headers := map[string][]string{}
	for k, values := range protoReq.Headers {
		headers[k] = values.Values
	}

	return &CallResourceRequest{
		PluginContext:	f.PluginContext(protoReq.PluginContext),
		Path:	protoReq.Path,
		Method:	protoReq.Method,
		URL:	protoReq.Url,
		Headers:	headers,
		Body:	protoReq.Body,
	}
}

func (f ConvertFromProtobuf) CallResourceResponse(protoResp *pluginv2.CallResourceResponse) *CallResourceResponse {
	headers := map[string][]string{}
	for k, values := range protoResp.Headers {
		headers[k] = values.Values
	}

	return &CallResourceResponse{
		Status:	int(protoResp.Code),
		Body:	protoResp.Body,
		Headers:	headers,
	}
}

func (f ConvertFromProtobuf) CheckHealthRequest(protoReq *pluginv2.CheckHealthRequest) *CheckHealthRequest {
	if protoReq.Headers == nil {
		protoReq.Headers = map[string]string{}
	}

	return &CheckHealthRequest{
		PluginContext:	f.PluginContext(protoReq.PluginContext),
		Headers:	protoReq.Headers,
	}
}

func (f ConvertFromProtobuf) CheckHealthResponse(protoResp *pluginv2.CheckHealthResponse) *CheckHealthResult {
	status := HealthStatusUnknown
	switch protoResp.Status {
	case pluginv2.CheckHealthResponse_ERROR:
		status = HealthStatusError
	case pluginv2.CheckHealthResponse_OK:
		status = HealthStatusOk
	}

	return &CheckHealthResult{
		Status:	status,
		Message:	protoResp.Message,
		JSONDetails:	protoResp.JsonDetails,
	}
}

func (f ConvertFromProtobuf) PluginInfoGetReq(protoReq *pluginv2.PluginInfoGetReq) *PluginInfoGetReq {

	return &PluginInfoGetReq{
		PluginContext: f.PluginContext(protoReq.PluginContext),
	}
}

func (f ConvertFromProtobuf) PluginInfoGetRes(protoResp *pluginv2.PluginInfoGetRes) *PluginInfoGetRes {

	return &PluginInfoGetRes{
		PluginID:	protoResp.PluginID,
		PluginVersion:	protoResp.PluginVersion,
	}
}
