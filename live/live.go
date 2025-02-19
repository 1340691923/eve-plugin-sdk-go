package live

import (
	"context"
	"errors"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/goccy/go-json"
)

type LiveChannel interface {
	Pub2Channel(ctx context.Context, req []byte) (res map[string]interface{}, err error)
}

type Live struct {
	registerChannelMap map[string]LiveChannel
}

func NewLive(registerChannelMap map[string]LiveChannel) *Live {
	return &Live{registerChannelMap: registerChannelMap}
}

func (live *Live) Pub2Channel(ctx context.Context, req *backend.Pub2ChannelRequest) (*backend.Pub2ChannelResponse, error) {
	if len(live.registerChannelMap) == 0 {
		return &backend.Pub2ChannelResponse{Status: backend.PubStatusError}, errors.New("注册频道数为0")
	}
	if _, ok := live.registerChannelMap[req.Channel]; !ok {
		return &backend.Pub2ChannelResponse{Status: backend.PubStatusError}, errors.New("没有找到该频道")
	}
	duck := live.registerChannelMap[req.Channel]
	res, err := duck.Pub2Channel(ctx, req.Data)
	if err != nil {
		return &backend.Pub2ChannelResponse{Status: backend.PubStatusError}, err
	}
	resB, err := json.Marshal(res)
	if err != nil {
		return &backend.Pub2ChannelResponse{Status: backend.PubStatusError}, err
	}
	return &backend.Pub2ChannelResponse{Status: backend.PubStatusOk, JsonDetails: resB}, nil
}
