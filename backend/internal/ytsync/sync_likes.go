package ytsync

import (
	"github.com/rwlist/youtube/internal/global"
	"github.com/rwlist/youtube/internal/models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/youtube/v3"
)

type LikedSync struct {
	cursor CursorCreator
}

func NewLikedSync(cursor CursorCreator) *LikedSync {
	return &LikedSync{
		cursor: cursor,
	}
}

func (s *LikedSync) Sync(syncID string, engine global.LikedEngine) error {
	log.WithField("id", syncID).Info("starting LikedSync.Sync")

	entry, err := engine.CatalogEntry()
	if err != nil {
		return err
	}

	opts := &CursorOpts{
		UserID:   entry.UserID,
		PageSize: 50,
	}
	cursor, err := s.cursor.CreateCursor(opts)
	if err != nil {
		return err
	}

	xord := ""
	for pageNum := 0; ; pageNum++ {
		resp, err := cursor.Next()
		if err != nil {
			return err
		}

		log.WithField("id", syncID).WithField("page", pageNum).Info("got page for likes sync")

		finished, err := s.updateSome(syncID, engine, resp.Items, &xord)
		if err != nil {
			return err
		}
		if finished {
			break
		}

		if int64(len(resp.Items)) < resp.PageInfo.ResultsPerPage {
			break
		}
	}

	log.WithField("id", syncID).Info("finished LikedSync.Sync")
	return nil
}

func (s *LikedSync) updateSome(_ string, engine global.LikedEngine, items []*youtube.PlaylistItem, xord *string) (finished bool, err error) {
	var ids []string
	for _, item := range items {
		ids = append(ids, item.ContentDetails.VideoId)
	}

	local, err := engine.FindByYoutubeIdsTyped(ids)
	if err != nil {
		return false, err
	}

	localMap := make(map[string]models.LikedModel)
	for i := range local {
		localMap[local[i].YoutubeID] = local[i]
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
			newItem, err := engine.InsertAfterTyped(*xord, models.LikedDataFromYoutube(data))
			if err != nil {
				return false, err
			}
			*xord = newItem.Xord
			continue
		}

		if data.PublishedAt != nil && (localItem.PublishedAt == nil || data.PublishedAt.After(*localItem.PublishedAt)) {
			finished = false
		}

		newItem, err := engine.UpdateDataTyped(models.LikedDataFromYoutube(data))
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
