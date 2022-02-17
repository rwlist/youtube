package global

import (
	"github.com/rwlist/youtube/internal/models"
)

type LikedSync interface {
	Sync(id string, engine LikedEngine) error
}

type CatalogLists interface {
	Update(list *models.CatalogList) error
}

type Directory struct {
	LikedSync    LikedSync
	CatalogsRepo CatalogLists
}
