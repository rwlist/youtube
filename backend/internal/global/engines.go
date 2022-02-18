package global

import "github.com/rwlist/youtube/internal/models"

type LikedEngine interface {
	CatalogEntry() (*models.CatalogList, error)
	FindTyped(youtubeIDs []string) ([]models.LikedModel, error)
	InsertAfter(xord string, data *models.LikedData) (*models.LikedModel, error)
	UpdateData(data *models.LikedData) (*models.LikedModel, error)
	MoveAfter(xord string, itemID uint) (newXord string, err error)
}
