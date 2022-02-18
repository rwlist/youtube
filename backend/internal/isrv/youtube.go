package isrv

import (
	"context"
	"github.com/rwlist/youtube/internal/proto"
	"github.com/rwlist/youtube/internal/rpc"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Youtube struct {
	oauthConfig *oauth2.Config
}

func NewYoutube(oauthConfig *oauth2.Config) *Youtube {
	return &Youtube{oauthConfig: oauthConfig}
}

func (y *Youtube) api(ctx context.Context) (*youtube.Service, error) {
	user := rpc.GetUser(ctx)
	cli := y.oauthConfig.Client(context.Background(), user.GoogleToken)
	api, err := youtube.NewService(context.Background(), option.WithHTTPClient(cli))
	return api, err
}

func (y *Youtube) Playlists(ctx context.Context) (*proto.Playlists, error) {
	api, err := y.api(ctx)
	if err != nil {
		return nil, err
	}

	call := api.Playlists.List([]string{"contentDetails", "id", "localizations", "player", "snippet", "status"})
	call.Id("LL", "WL")
	//call.Mine(true)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	return &proto.Playlists{
		Response: resp,
	}, nil
}

func (y *Youtube) Liked(ctx context.Context) (*proto.PlaylistItems, error) {
	api, err := y.api(ctx)
	if err != nil {
		return nil, err
	}

	call := api.PlaylistItems.List([]string{"contentDetails", "id", "snippet", "status"})
	call.PlaylistId("LL")
	call.MaxResults(50)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	return &proto.PlaylistItems{
		Response: resp,
	}, nil
}
