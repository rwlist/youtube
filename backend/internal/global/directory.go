package global

import (
	"github.com/rwlist/youtube/internal/models"
	"gorm.io/gorm"
)

type LikedSync interface {
	Sync(id string, engine LikedEngine) error
}

type CatalogLists interface {
	Update(tx *gorm.DB, list *models.CatalogList) error
}

type Directory struct {
	LikedSync    LikedSync
	CatalogsRepo CatalogLists
}
