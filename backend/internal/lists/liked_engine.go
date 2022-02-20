package lists

import (
	"time"

	"github.com/rwlist/youtube/internal/global"
	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/proto"
	log "github.com/sirupsen/logrus"
)

type LikedEngine struct {
	CommonEngine
}

func NewLikedEngine(s *Storage, dir *global.Directory) *LikedEngine {
	return &LikedEngine{
		CommonEngine: NewCommonEngine(
			s,
			dir,
			&CommonOpts{
				newModel: func() interface{} { return &models.LikedModel{} },
				newArray: func() interface{} { return &[]models.LikedModel{} },
				convertArray: func(i interface{}) ([]proto.ListItem, error) {
					arr, ok := i.(*[]models.LikedModel)
					if !ok {
						return nil, ErrInvalidType
					}
					var items []proto.ListItem
					for i := range *arr {
						item := (*arr)[i]
						items = append(items, item)
					}
					return items, nil
				},
				initMeta: models.CatalogMeta{
					ObjectIDField:    "youtube_id",
					IsUniqueObjectID: true,
				},
			},
		),
	}
}

func (e LikedEngine) Transaction(f func(e Engine) error) error {
	return e.storage.Transaction(func(s *Storage) error {
		e.storage = s
		return f(&e)
	})
}

func (e *LikedEngine) ListItemsTyped() ([]models.LikedModel, error) {
	var items []models.LikedModel
	err := e.storage.FindAll(&items)
	return items, err
}

func (e *LikedEngine) FindByYoutubeIdsTyped(youtubeIDs []string) ([]models.LikedModel, error) {
	var items []models.LikedModel
	err := e.storage.FindByObjectIDs(&items, youtubeIDs)
	return items, err
}

func (e *LikedEngine) InsertAfterTyped(xord string, data *models.LikedData) (*models.LikedModel, error) {
	newXord, err := e.xordForInsert(xord)
	if err != nil {
		return nil, err
	}

	item := models.LikedModel{
		Model:     models.XordModel(newXord),
		LikedData: *data,
	}
	err = e.storage.Insert(&item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (e *LikedEngine) UpdateDataTyped(data *models.LikedData) (*models.LikedModel, error) {
	var model models.LikedModel
	// TODO: don't fetch full model to update
	err := e.storage.FirstByKey(&model, e.storage.Catalog().Meta.ObjectIDField, data.YoutubeID)
	if err != nil {
		return nil, err
	}

	model.LikedData = *data

	err = e.storage.Save(&model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (e *LikedEngine) ExecuteQuery(query string) (*proto.QueryResponse, error) {
	switch query {
	case ":sync":
		return e.StartSync()
	}

	return &proto.QueryResponse{
		Status: "unknown query",
	}, nil
}

func (e *LikedEngine) StartSync() (*proto.QueryResponse, error) {
	id := time.Now().String()

	ch := make(chan struct{})

	go func() {
		defer close(ch)

		err := e.Transaction(func(e Engine) error {
			engine := e.(*LikedEngine)
			return engine.dir.LikedSync.Sync(id, engine)
		})
		if err != nil {
			log.WithField("id", id).WithError(err).Error("sync tx failed")
		}
	}()

	return &proto.QueryResponse{
		Status:      "Sync is started, id=" + id,
		DoneChannel: ch,
	}, nil
}
