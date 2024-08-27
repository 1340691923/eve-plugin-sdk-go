package backend

import (
	"context"
	"net/http"
	"strconv"
)

type CheckHealthHandler interface {
	CheckHealth(ctx context.Context, req *CheckHealthRequest) (*CheckHealthResult, error)
}

type CheckHealthHandlerFunc func(ctx context.Context, req *CheckHealthRequest) (*CheckHealthResult, error)

func (fn CheckHealthHandlerFunc) CheckHealth(ctx context.Context, req *CheckHealthRequest) (*CheckHealthResult, error) {
	return fn(ctx, req)
}

type HealthStatus int

const (
	HealthStatusUnknown	HealthStatus	= iota

	HealthStatusOk

	HealthStatusError
)

var healthStatusNames = map[int]string{
	0:	"UNKNOWN",
	1:	"OK",
	2:	"ERROR",
}

func (hs HealthStatus) String() string {
	s, exists := healthStatusNames[int(hs)]
	if exists {
		return s
	}
	return strconv.Itoa(int(hs))
}

type CheckHealthRequest struct {
	PluginContext	PluginContext

	Headers	map[string]string
}

func (req *CheckHealthRequest) SetHTTPHeader(key, value string) {
	if req.Headers == nil {
		req.Headers = map[string]string{}
	}

	setHTTPHeaderInStringMap(req.Headers, key, value)
}

func (req *CheckHealthRequest) DeleteHTTPHeader(key string) {
	deleteHTTPHeaderInStringMap(req.Headers, key)
}

func (req *CheckHealthRequest) GetHTTPHeader(key string) string {
	return req.GetHTTPHeaders().Get(key)
}

func (req *CheckHealthRequest) GetHTTPHeaders() http.Header {
	return getHTTPHeadersFromStringMap(req.Headers)
}

type CheckHealthResult struct {
	Status	HealthStatus

	Message	string

	JSONDetails	[]byte
}

type PluginInfoHandler interface {
	PluginInfo(ctx context.Context, req *PluginInfoGetReq) (*PluginInfoGetRes, error)
}

type PluginInfoHandlerFunc func(ctx context.Context, req *PluginInfoGetReq) (*PluginInfoGetRes, error)

func (fn PluginInfoHandlerFunc) PluginInfo(ctx context.Context, req *PluginInfoGetReq) (*PluginInfoGetRes, error) {
	return fn(ctx, req)
}

type PluginInfoGetReq struct {
	PluginContext PluginContext
}

type PluginInfoGetRes struct {
	PluginID	string
	PluginVersion	string
}

var _ ForwardHTTPHeaders = (*CheckHealthRequest)(nil)
