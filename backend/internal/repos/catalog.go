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
	if err != nil {
		return nil, err
	}
	for i := range lists {
		err := lists[i].AfterLoad()
		if err != nil {
			return nil, err
		}
	}
	return lists, nil
}

// Create will generate TableName automatically.
func (r *CatalogLists) Create(list *models.CatalogList, txCallback func(*gorm.DB) error) error {
	err := list.GenerateTableName()
	if err != nil {
		return err
	}

	err = list.BeforeSave()
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

func (r *CatalogLists) Update(list *models.CatalogList) error {
	err := list.BeforeSave()
	if err != nil {
		return err
	}

	return r.db.Save(list).Error
}

func (r *CatalogLists) GetByListID(userID uint, listID string) (*models.CatalogList, error) {
	var list models.CatalogList
	err := r.db.Where("user_id = ? AND list_id = ?", userID, listID).First(&list).Error
	if err != nil {
		return nil, err
	}
	err = list.AfterLoad()
	if err != nil {
		return nil, err
	}
	return &list, nil
}
