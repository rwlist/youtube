package rpc

import (
	"github.com/rwlist/gjrpc/pkg/jsonrpc"
	"github.com/rwlist/youtube/internal/proto"
	"strings"
)

func ConvertError(err error) (jsonrpc.Result, *jsonrpc.Error) {
	if strings.Contains(err.Error(), "oauth2: token expired and refresh token is not set") {
		return nil, proto.AuthError.WithData("oauth2: token expired and refresh token is not set")
	}

	return nil, proto.InternalError.WithData(err)
}
