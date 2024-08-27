package backend

import (
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
	"time"
)

type ConvertToProtobuf struct{}

func ToProto() ConvertToProtobuf {
	return ConvertToProtobuf{}
}

func (t ConvertToProtobuf) User(user *User) *pluginv2.User {
	if user == nil {
		return nil
	}

	return &pluginv2.User{
		Login:	user.Login,
		Name:	user.Name,
		Email:	user.Email,
		Role:	user.Role,
	}
}

func (t ConvertToProtobuf) AppInstanceSettings(s *AppInstanceSettings) *pluginv2.AppInstanceSettings {
	if s == nil {
		return nil
	}

	return &pluginv2.AppInstanceSettings{
		JsonData:	s.JSONData,
		DecryptedSecureJsonData:	s.DecryptedSecureJSONData,
		LastUpdatedMS:	s.Updated.UnixNano() / int64(time.Millisecond),
	}
}

func (t ConvertToProtobuf) DataSourceInstanceSettings(s *DataSourceInstanceSettings) *pluginv2.DataSourceInstanceSettings {
	if s == nil {
		return nil
	}

	return &pluginv2.DataSourceInstanceSettings{
		Id:	s.ID,
		Uid:	s.UID,
		Name:	s.Name,
		Url:	s.URL,
		User:	s.User,
		Database:	s.Database,
		BasicAuthEnabled:	s.BasicAuthEnabled,
		BasicAuthUser:	s.BasicAuthUser,
		JsonData:	s.JSONData,
		DecryptedSecureJsonData:	s.DecryptedSecureJSONData,
		LastUpdatedMS:	s.Updated.UnixNano() / int64(time.Millisecond),
	}
}

func (t ConvertToProtobuf) PluginContext(pluginCtx PluginContext) *pluginv2.PluginContext {
	return &pluginv2.PluginContext{
		OrgId:	pluginCtx.OrgID,
		PluginId:	pluginCtx.PluginID,
		User:	t.User(pluginCtx.User),
		AppInstanceSettings:	t.AppInstanceSettings(pluginCtx.AppInstanceSettings),
		DataSourceInstanceSettings:	t.DataSourceInstanceSettings(pluginCtx.DataSourceInstanceSettings),
	}
}

func (t ConvertToProtobuf) CallResourceResponse(resp *CallResourceResponse) *pluginv2.CallResourceResponse {
	headers := map[string]*pluginv2.StringList{}

	for key, values := range resp.Headers {
		headers[key] = &pluginv2.StringList{Values: values}
	}

	return &pluginv2.CallResourceResponse{
		Headers:	headers,
		Code:	int32(resp.Status),
		Body:	resp.Body,
	}
}

func (t ConvertToProtobuf) CallResourceRequest(req *CallResourceRequest) *pluginv2.CallResourceRequest {
	protoReq := &pluginv2.CallResourceRequest{
		PluginContext:	t.PluginContext(req.PluginContext),
		Path:	req.Path,
		Method:	req.Method,
		Url:	req.URL,
		Body:	req.Body,
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

func (f ConvertToProtobuf) PluginInfoGetRes(protoResp *PluginInfoGetRes) *pluginv2.PluginInfoGetRes {

	return &pluginv2.PluginInfoGetRes{
		PluginID:	protoResp.PluginID,
		PluginVersion:	protoResp.PluginVersion,
	}
}

func (t ConvertToProtobuf) CheckHealthResponse(res *CheckHealthResult) *pluginv2.CheckHealthResponse {
	return &pluginv2.CheckHealthResponse{
		Status:	t.HealthStatus(res.Status),
		Message:	res.Message,
		JsonDetails:	res.JSONDetails,
	}
}

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
