package backend

import (
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
	"net/http"
)

type resourceSDKAdapter struct {
	callResourceHandler CallResourceHandler
}

func newResourceSDKAdapter(handler CallResourceHandler) *resourceSDKAdapter {
	return &resourceSDKAdapter{
		callResourceHandler: handler,
	}
}

type callResourceResponseSenderFunc func(resp *CallResourceResponse) error

func (fn callResourceResponseSenderFunc) Send(resp *CallResourceResponse) error {
	return fn(resp)
}

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
