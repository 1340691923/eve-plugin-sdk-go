package backend

import (
	"context"
	"strconv"
)

type PubStatus int

const (
	PubStatusUnknown PubStatus = iota

	PubStatusOk

	PubStatusError
)

var pubStatusNames = map[int]string{
	0: "UNKNOWN",
	1: "OK",
	2: "ERROR",
}

type Pub2ChannelRequest struct {
	PluginContext PluginContext
	Channel       string
	Data          []byte
}

type Pub2ChannelResponse struct {
	Status      PubStatus
	Message     string
	JsonDetails []byte
}

func (hs PubStatus) String() string {
	s, exists := pubStatusNames[int(hs)]
	if exists {
		return s
	}
	return strconv.Itoa(int(hs))
}

type LiveHandler interface {
	Pub2Channel(ctx context.Context, req *Pub2ChannelRequest) (*Pub2ChannelResponse, error)
}

type LiveHandlerFunc func(ctx context.Context, req *Pub2ChannelRequest) (*Pub2ChannelResponse, error)

func (fn LiveHandlerFunc) Pub2Channel(ctx context.Context, req *Pub2ChannelRequest) (*Pub2ChannelResponse, error) {
	return fn(ctx, req)
}
