package lists

import (
	"github.com/rwlist/youtube/internal/global"
	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/proto"
	log "github.com/sirupsen/logrus"
	"time"
)

type Engine interface {
	InitStorage() error
	Info() (*proto.ListInfo, error)
	ListItems() ([]proto.ListItem, error)
	StartSync() (proto.ListSync, error)
	PageItems(req proto.PageRequest) ([]proto.ListItem, error)
}

type LikedEngine struct {
	s   *Storage
	dir *global.Directory
}

func NewLikedEngine(s *Storage, dir *global.Directory) *LikedEngine {
	return &LikedEngine{
		s:   s,
		dir: dir,
	}
}

func (e *LikedEngine) InitStorage() error {
	return e.s.CreateTable(&models.ListDataUnique{})
}

func (e *LikedEngine) Automigrate() error {
	return e.s.Automigrate(&models.ListDataUnique{})
}

func (e LikedEngine) ReplaceStorage(s *Storage) *LikedEngine {
	e.s = s
	return &e
}

func (e *LikedEngine) CatalogEntry() (*models.CatalogList, error) {
	return e.s.Catalog(), nil
}

func (e *LikedEngine) Info() (*proto.ListInfo, error) {
	info := e.s.Catalog().ToInfo()
	itemsCount, err := e.s.CountAll()
	if err != nil {
		return nil, err
	}
	info.ItemsCount = itemsCount
	return info, nil
}

func convertItem(item models.ListDataUnique) proto.ListItem {
	return proto.ListItem{
		YoutubeID: item.YoutubeID,
		Title:     item.Title,
		Author:    item.Author,
		ChannelID: item.ChannelID,
		ItemID:    item.ItemID,
		Xord:      item.Xord,
	}
}

func (e *LikedEngine) ListItems() ([]proto.ListItem, error) {
	var items []models.ListDataUnique
	err := e.s.FindAll(&items)
	if err != nil {
		return nil, err
	}

	res := make([]proto.ListItem, len(items))
	for i, item := range items {
		res[i] = convertItem(item)
	}
	return res, nil
}

func (e *LikedEngine) PageItems(req proto.PageRequest) ([]proto.ListItem, error) {
	var items []models.ListDataUnique
	err := e.s.FindByPageRequest(req, &items)
	if err != nil {
		return nil, err
	}

	res := make([]proto.ListItem, len(items))
	for i, item := range items {
		res[i] = convertItem(item)
	}
	return res, nil
}

func (e *LikedEngine) FindTyped(youtubeIDs []string) ([]models.ListDataUnique, error) {
	var items []models.ListDataUnique
	err := e.s.FindByYoutubeID(&items, youtubeIDs)
	return items, err
}

func (e *LikedEngine) InsertAfter(xord string, data *models.YoutubeData) (*models.ListDataUnique, error) {
	newXord, err := e.prepareXord(xord)
	if err != nil {
		return nil, err
	}

	item := models.ListDataUnique{
		Xord:        newXord,
		YoutubeID:   data.YoutubeID,
		Title:       data.Title,
		Author:      data.Author,
		ChannelID:   data.ChannelID,
		PublishedAt: data.PublishedAt,
	}
	err = e.s.Insert(&item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// "" is beginning
func (e *LikedEngine) prepareXord(at string) (string, error) {
	cnt := 2
	if at != "" {
		cnt = 3
	}

	var start []models.ListDataUnique
	err := e.s.OrderLimit(&start, at, cnt)
	if err != nil {
		return "", err
	}

	// <at> <new> <[0]> <[1]>
	newXord := ""
	if len(start) == 0 {
		newXord = splitXord("", "")
	} else if at != "" && len(start) == 1 {
		newXord = splitXord(start[0].Xord, "")
	} else {
		if at != "" {
			start = start[1:]
		}
		newXord = start[0].Xord
		err = e.shiftRight(start)
		if err != nil {
			return "", err
		}
	}

	return newXord, nil
}

func (e *LikedEngine) UpdateData(data *models.YoutubeData) (*models.ListDataUnique, error) {
	var model models.ListDataUnique
	// TODO: don't fetch full model to update
	err := e.s.FirstByYoutubeID(&model, data.YoutubeID)
	if err != nil {
		return nil, err
	}

	model.UpdateData(data)
	err = e.s.Save(&model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (e *LikedEngine) MoveAfter(xord string, itemID uint) (newXord string, err error) {
	newXord, err = e.prepareXord(xord)
	if err != nil {
		return "", err
	}

	return newXord, e.s.UpdateXord(itemID, newXord)
}

func (e *LikedEngine) shiftRight(nxt2 []models.ListDataUnique) error {
	if len(nxt2) == 0 {
		return nil
	}
	l := nxt2[0].Xord
	r := ""
	if len(nxt2) == 2 {
		r = nxt2[1].Xord
	}

	newXord := splitXord(l, r)
	return e.s.UpdateXord(nxt2[0].ItemID, newXord)
}

func (e LikedEngine) Transaction(f func(e Engine) error) error {
	return e.s.Transaction(func(s *Storage) error {
		return f(e.ReplaceStorage(s))
	})
}

func (e *LikedEngine) StartSync() (proto.ListSync, error) {
	id := time.Now().String()

	go func() {
		err := e.Transaction(func(e Engine) error {
			engine := e.(*LikedEngine)
			return engine.dir.LikedSync.Sync(id, engine)
		})
		if err != nil {
			log.WithField("id", id).WithError(err).Error("sync tx failed")
		}
	}()

	return proto.ListSync{
		Status: "Sync is started, id=" + id,
	}, nil
}
