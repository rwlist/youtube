package global

import "github.com/rwlist/youtube/internal/models"

type LikedEngine interface {
	CatalogEntry() (*models.CatalogList, error)
	FindTyped(youtubeIDs []string) ([]models.ListDataUnique, error)
	InsertAfter(xord string, data *models.YoutubeData) (*models.ListDataUnique, error)
	UpdateData(data *models.YoutubeData) (*models.ListDataUnique, error)
	MoveAfter(xord string, itemID uint) (newXord string, err error)
}
