package isrv

import (
	"context"
	"fmt"
	"github.com/rwlist/youtube/internal/lists"
	"github.com/rwlist/youtube/internal/proto"
	"github.com/rwlist/youtube/internal/rpc"
	"golang.org/x/oauth2"
)

type Lists struct {
	oauthConfig *oauth2.Config
	youtube     *Youtube
	catalog     *lists.Catalog
}

func NewLists(oauthConfig *oauth2.Config, youtube *Youtube, catalog *lists.Catalog) *Lists {
	return &Lists{
		oauthConfig: oauthConfig,
		youtube:     youtube,
		catalog:     catalog,
	}
}

func (l *Lists) All(ctx context.Context) (proto.AllLists, error) {
	user := rpc.GetUser(ctx)

	all, err := l.catalog.UserLists(user.ID)
	if err != nil {
		return proto.AllLists{}, err
	}

	return proto.AllLists{
		Lists: all,
	}, nil
}

// TODO: write better
func (l *Lists) ListInfo(ctx context.Context, listID string) (proto.ListInfo, error) {
	all, err := l.All(ctx)
	if err != nil {
		return proto.ListInfo{}, err
	}

	for _, list := range all.Lists {
		if list.ID == listID {
			return list, nil
		}
	}

	return proto.ListInfo{}, fmt.Errorf("list not found")
}

func (l *Lists) ListItems(ctx context.Context, listID string) (proto.ListItems, error) {
	user := rpc.GetUser(ctx)

	items, err := l.catalog.ViewList(user.ID, listID)
	if err != nil {
		return proto.ListItems{}, err
	}

	return proto.ListItems{
		Items: items,
	}, nil
}
