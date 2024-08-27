package backend

import (
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
)

type pluginInfoSDKAdapter struct {
	pluginInfoHandler	PluginInfoHandler
	checkHealthHandler	CheckHealthHandler
}

func newPluginInfoSDKAdapter(pluginInfoHandler PluginInfoHandler, checkHealthHandler CheckHealthHandler) *pluginInfoSDKAdapter {
	return &pluginInfoSDKAdapter{
		pluginInfoHandler:	pluginInfoHandler,
		checkHealthHandler:	checkHealthHandler,
	}
}

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
		PluginID:	"",
		PluginVersion:	"",
	}, nil
}

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
