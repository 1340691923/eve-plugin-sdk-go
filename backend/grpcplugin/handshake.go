package grpcplugin

import (
	plugin "github.com/hashicorp/go-plugin"
)

const (
	ProtocolVersion	= 2

	MagicCookieKey	= "ev_author"

	MagicCookieValue	= "xiaowenlong"
)

var handshake = plugin.HandshakeConfig{
	ProtocolVersion:	ProtocolVersion,
	MagicCookieKey:	MagicCookieKey,
	MagicCookieValue:	MagicCookieValue,
}
