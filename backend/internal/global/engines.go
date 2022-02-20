package global

import "github.com/rwlist/youtube/internal/models"

type LikedEngine interface {
	CatalogEntry() (*models.CatalogList, error)
	FindByYoutubeIdsTyped(youtubeIDs []string) ([]models.LikedModel, error)
	InsertAfterTyped(xord string, data *models.LikedData) (*models.LikedModel, error)
	UpdateDataTyped(data *models.LikedData) (*models.LikedModel, error)
	MoveAfter(xord string, itemID uint) (newXord string, err error)
}
