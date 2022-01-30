package isrv

import (
	"context"
	"fmt"
	"github.com/rwlist/youtube/internal/proto"
	"golang.org/x/oauth2"
)

// TODO: implement real logic later

type Lists struct {
	oauthConfig *oauth2.Config
}

func NewLists(oauthConfig *oauth2.Config) *Lists {
	return &Lists{oauthConfig: oauthConfig}
}

func (l *Lists) All(ctx context.Context) (proto.AllLists, error) {
	return proto.AllLists{
		Lists: []proto.ListInfo{
			{
				ID:       "liked",
				Name:     "Liked videos",
				ListType: "external",
			},
			{
				ID:       "history",
				Name:     "History",
				ListType: "external",
			},
			{
				ID:       "view-later",
				Name:     "View later",
				ListType: "external",
			},
		},
	}, nil
}

func (l *Lists) ListInfo(ctx context.Context, listID string) (proto.ListInfo, error) {
	lists, err := l.All(ctx)
	if err != nil {
		return proto.ListInfo{}, err
	}

	for _, list := range lists.Lists {
		if list.ID == listID {
			return list, nil
		}
	}

	return proto.ListInfo{}, fmt.Errorf("list not found")
}

func (l *Lists) ListItems(ctx context.Context, listID string) (proto.ListItems, error) {
	if listID != "history" {
		return proto.ListItems{}, fmt.Errorf("list not found")
	}

	return proto.ListItems{
		Items: []proto.ListItem{
			{
				YoutubeID: "5ESJH1NLMLs",
				Title:     "Children Of The Magenta Line",
				Author:    "Mossie Fly",
				ChannelID: "UCGIkFNbztHRaX0GB78SWaZA",
				Xord:      "aba",
				ItemID:    "0",
			},
			{
				YoutubeID: "68T9EFlCsUc",
				Title:     "Making Music Is Easy",
				Author:    "Eliminate",
				ChannelID: "UCI7kKmUuSQOHUvSWIYFDf1Q",
				Xord:      "acc",
				ItemID:    "1",
			},
		},
	}, nil
}
