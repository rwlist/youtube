package isrv

import (
	"context"
	"github.com/rwlist/youtube/internal/lists"
	"github.com/rwlist/youtube/internal/proto"
	"github.com/rwlist/youtube/internal/rpc"
	"golang.org/x/oauth2"
)

type List struct {
	oauthConfig *oauth2.Config
	catalog     *lists.Catalog
}

func NewList(oauthConfig *oauth2.Config, catalog *lists.Catalog) *List {
	return &List{
		oauthConfig: oauthConfig,
		catalog:     catalog,
	}
}

func (l *List) Info(ctx context.Context, listID string) (proto.ListInfo, error) {
	user := rpc.GetUser(ctx)

	engine, err := l.catalog.GetList(user.ID, listID)
	if err != nil {
		return proto.ListInfo{}, err
	}

	info, err := engine.Info()
	if err != nil {
		return proto.ListInfo{}, err
	}

	return *info, nil
}

func (l *List) Items(ctx context.Context, listID string) (proto.ListItems, error) {
	user := rpc.GetUser(ctx)

	engine, err := l.catalog.GetList(user.ID, listID)
	if err != nil {
		return proto.ListItems{}, err
	}

	items, err := engine.ListItems()
	if err != nil {
		return proto.ListItems{}, err
	}

	return proto.ListItems{
		Items: items,
	}, nil
}

func (l *List) Sync(ctx context.Context, listID string) (proto.ListSync, error) {
	user := rpc.GetUser(ctx)

	engine, err := l.catalog.GetList(user.ID, listID)
	if err != nil {
		return proto.ListSync{}, err
	}

	return engine.StartSync()
}
