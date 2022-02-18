package ytsync

import (
	"context"
	"github.com/rwlist/youtube/internal/global"
	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/repos"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type LikedSync struct {
	oauthConfig *oauth2.Config
	users       *repos.Users
}

func NewLikedSync(oauthConfig *oauth2.Config, users *repos.Users) *LikedSync {
	return &LikedSync{
		oauthConfig: oauthConfig,
		users:       users,
	}
}

func (s *LikedSync) Sync(id string, engine global.LikedEngine) error {
	log.WithField("id", id).Info("starting LikedSync.Sync")

	entry, err := engine.CatalogEntry()
	if err != nil {
		return err
	}

	userID := entry.UserID
	user, err := s.users.Get(userID)
	if err != nil {
		return err
	}

	err = user.PreloadToken()
	if err != nil {
		return err
	}

	cli := s.oauthConfig.Client(context.Background(), user.GoogleToken)
	api, err := youtube.NewService(context.Background(), option.WithHTTPClient(cli))
	if err != nil {
		return err
	}

	call := api.PlaylistItems.List([]string{"contentDetails", "id", "snippet", "status"})
	call.PlaylistId("LL")
	call.MaxResults(50)

	xord := ""

	for pageNum := 0; ; pageNum++ {
		resp, err := call.Do()
		if err != nil {
			return err
		}

		log.WithField("id", id).WithField("page", pageNum).Info("got page for likes sync")

		finished, err := s.updateSome(id, engine, resp.Items, &xord)
		if err != nil {
			return err
		}
		if finished {
			break
		}

		if int64(len(resp.Items)) < resp.PageInfo.ResultsPerPage {
			break
		}
		call.PageToken(resp.NextPageToken)
	}

	log.WithField("id", id).Info("finished LikedSync.Sync")
	return nil
}

func (s *LikedSync) updateSome(id string, engine global.LikedEngine, items []*youtube.PlaylistItem, xord *string) (finished bool, err error) {
	var ids []string
	for _, item := range items {
		ids = append(ids, item.ContentDetails.VideoId)
	}

	local, err := engine.FindTyped(ids)
	if err != nil {
		return false, err
	}

	localMap := make(map[string]models.LikedModel)
	for _, l := range local {
		localMap[l.YoutubeID] = l
	}

	// assume that nothing is changed yet
	finished = true

	for _, item := range items {
		data, err := models.FromYoutubeItem(item)
		if err != nil {
			return false, err
		}

		localItem, ok := localMap[data.YoutubeID]
		if !ok {
			finished = false
			newItem, err := engine.InsertAfter(*xord, models.LikedDataFromYoutube(data))
			if err != nil {
				return false, err
			}
			*xord = newItem.Xord
			continue
		}

		if data.PublishedAt != nil && (localItem.PublishedAt == nil || data.PublishedAt.After(*localItem.PublishedAt)) {
			finished = false
		}

		newItem, err := engine.UpdateData(models.LikedDataFromYoutube(data))
		if err != nil {
			return false, err
		}

		newXord, err := engine.MoveAfter(*xord, newItem.ItemID)
		if err != nil {
			return false, err
		}
		*xord = newXord
	}

	return finished, nil
}
