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

	//gjrpc:handle-route proto.ListService
	Lists ListsImpl
}

func NewHandlers(auth AuthImpl, youtube YoutubeImpl, lists ListsImpl) Handlers {
	return Handlers{
		Auth:    auth,
		Youtube: youtube,
		Lists:   lists,
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

type ListsImpl interface {
	All(ctx context.Context) (proto.AllLists, error)
	ListInfo(ctx context.Context, listID string) (proto.ListInfo, error)
	ListItems(ctx context.Context, listID string) (proto.ListItems, error)
}
