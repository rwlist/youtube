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

	//gjrpc:handle-route proto.ListsCatalog
	Catalog CatalogImpl

	//gjrpc:handle-route proto.ListService
	List ListImpl
}

func NewHandlers(auth AuthImpl, youtube YoutubeImpl, catalog CatalogImpl, list ListImpl) *Handlers {
	return &Handlers{
		Auth:    auth,
		Youtube: youtube,
		Catalog: catalog,
		List:    list,
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

type CatalogImpl interface {
	All(ctx context.Context) (proto.AllLists, error)
}

type ListImpl interface {
	Info(ctx context.Context, listID string) (proto.ListInfo, error)
	Items(ctx context.Context, listID string) (proto.ListItems, error)
	PageItems(ctx context.Context, req proto.PageRequest) (proto.ListItems, error)
	ExecuteQuery(ctx context.Context, query proto.Query) (proto.QueryResponse, error)
}
