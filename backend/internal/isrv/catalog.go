package isrv

import (
	"context"

	"github.com/rwlist/youtube/internal/lists"
	"github.com/rwlist/youtube/internal/proto"
	"github.com/rwlist/youtube/internal/rpc"
	"golang.org/x/oauth2"
)

type ListsCatalog struct {
	oauthConfig *oauth2.Config
	catalog     *lists.Catalog
}

func NewListsCatalog(oauthConfig *oauth2.Config, catalog *lists.Catalog) *ListsCatalog {
	return &ListsCatalog{
		oauthConfig: oauthConfig,
		catalog:     catalog,
	}
}

func (l *ListsCatalog) All(ctx context.Context) (*proto.AllLists, error) {
	user := rpc.GetUser(ctx)

	engines, err := l.catalog.UserLists(user.ID)
	if err != nil {
		return nil, err
	}

	var infos []proto.ListInfo
	for _, engine := range engines {
		info, err := engine.Info()
		if err != nil {
			return nil, err
		}
		infos = append(infos, *info)
	}

	return &proto.AllLists{
		Lists: infos,
	}, nil
}
