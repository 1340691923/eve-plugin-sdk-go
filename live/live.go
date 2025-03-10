package live

import (
	"context"
	"errors"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/goccy/go-json"
)

type LiveHandle interface {
	Pub2Channel(ctx context.Context, channel string, req []byte) (res map[string]interface{}, err error)
}

type Live struct {
	liveHandle LiveHandle
}

func NewLive(liveHandle LiveHandle) *Live {
	return &Live{liveHandle: liveHandle}
}

func (this *Live) Pub2Channel(ctx context.Context, req *backend.Pub2ChannelRequest) (*backend.Pub2ChannelResponse, error) {
	if this.liveHandle == nil {
		return &backend.Pub2ChannelResponse{Status: backend.PubStatusError}, errors.New("该插件没有实现长连接处理器")
	}

	res, err := this.liveHandle.Pub2Channel(ctx, req.Channel, req.Data)
	if err != nil {
		return &backend.Pub2ChannelResponse{Status: backend.PubStatusError}, err
	}
	resB, err := json.Marshal(res)
	if err != nil {
		return &backend.Pub2ChannelResponse{Status: backend.PubStatusError}, err
	}
	return &backend.Pub2ChannelResponse{Status: backend.PubStatusOk, JsonDetails: resB}, nil
}
