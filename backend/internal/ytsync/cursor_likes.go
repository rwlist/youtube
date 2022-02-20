package ytsync

import (
	"context"

	"github.com/rwlist/youtube/internal/repos"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type LikesCursor struct {
	oauthConfig *oauth2.Config
	users       *repos.Users

	call *youtube.PlaylistItemsListCall
}

func NewLikesCursor(oauthConfig *oauth2.Config, users *repos.Users) LikesCursor {
	return LikesCursor{
		oauthConfig: oauthConfig,
		users:       users,
	}
}

func (l LikesCursor) CreateCursor(opts *CursorOpts) (Cursor, error) {
	userID := opts.UserID
	user, err := l.users.Get(userID)
	if err != nil {
		return nil, err
	}

	err = user.PreloadToken()
	if err != nil {
		return nil, err
	}

	cli := l.oauthConfig.Client(context.Background(), user.GoogleToken)
	api, err := youtube.NewService(context.Background(), option.WithHTTPClient(cli))
	if err != nil {
		return nil, err
	}

	call := api.PlaylistItems.List([]string{"contentDetails", "id", "snippet", "status"})
	call.PlaylistId("LL")
	call.MaxResults(opts.PageSize)

	l.call = call
	return &l, nil
}

func (l *LikesCursor) Next() (*youtube.PlaylistItemListResponse, error) {
	resp, err := l.call.Do()
	if err != nil {
		return nil, err
	}

	l.call.PageToken(resp.NextPageToken)
	return resp, nil
}
