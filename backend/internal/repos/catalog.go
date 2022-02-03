package repos

import (
	"github.com/rwlist/youtube/internal/models"
	"gorm.io/gorm"
)

type CatalogLists struct {
	db *gorm.DB
}

func NewCatalogLists(db *gorm.DB) *CatalogLists {
	return &CatalogLists{db}
}

func (r *CatalogLists) FindUserLists(userID uint) ([]models.CatalogList, error) {
	var lists []models.CatalogList
	err := r.db.Where("user_id = ?", userID).Find(&lists).Error
	return lists, err
}

// Create will generate TableName automatically.
func (r *CatalogLists) Create(list *models.CatalogList, txCallback func(*gorm.DB) error) error {
	err := list.GenerateTableName()
	if err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(list).Error
		if err != nil {
			return err
		}
		return txCallback(tx)
	})
}

func (r *CatalogLists) GetByListID(userID uint, listID string) (*models.CatalogList, error) {
	var list models.CatalogList
	err := r.db.Where("user_id = ? AND list_id = ?", userID, listID).First(&list).Error
	return &list, err
}
