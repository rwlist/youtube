//go:generate gjrpc gen:server:router --protoPkg=../proto --handlersStruct=Handlers --out=gen_router.go
package rpc

import (
	"context"
	"github.com/rwlist/youtube/internal/proto"
)

type Handlers struct {
	//gjrpc:handle-route proto.Auth
	Auth AuthImpl

	//gjrpc:handle-route proto.Youtube
	Youtube YoutubeImpl
}

func NewHandlers(auth AuthImpl, youtube YoutubeImpl) Handlers {
	return Handlers{
		Auth:    auth,
		Youtube: youtube,
	}
}

type AuthImpl interface {
	Oauth(ctx context.Context) (proto.OAuthResponse, error)
	Status(ctx context.Context) (proto.AuthStatus, error)
}

type YoutubeImpl interface {
	Playlists(ctx context.Context) (proto.Playlists, error)
	Liked(ctx context.Context) (proto.PlaylistItems, error)
}
