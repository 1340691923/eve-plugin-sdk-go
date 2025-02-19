package backend

import (
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
)

type liveSDKAdapter struct {
	liveHandler LiveHandler
}

func newLiveSDKAdapter(handler LiveHandler) *liveSDKAdapter {
	return &liveSDKAdapter{
		liveHandler: handler,
	}
}

type liveResponseSenderFunc func(resp *Pub2ChannelResponse) error

func (fn liveResponseSenderFunc) Send(resp *Pub2ChannelResponse) error {
	return fn(resp)
}

func (a *liveSDKAdapter) Pub2Channel(ctx context.Context, protoReq *pluginv2.Pub2ChannelRequest) (*pluginv2.Pub2ChannelResponse, error) {
	if a.liveHandler != nil {
		parsedReq := FromProto().Pub2ChannelRequest(protoReq)

		res, err := a.liveHandler.Pub2Channel(ctx, parsedReq)
		if err != nil {
			return nil, err
		}
		return ToProto().Pub2ChannelResponse(res), nil
	}

	return &pluginv2.Pub2ChannelResponse{
		Status: pluginv2.Pub2ChannelResponse_OK,
	}, nil
}
