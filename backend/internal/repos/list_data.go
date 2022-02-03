package repos

import (
	"github.com/rwlist/youtube/internal/models"
	"gorm.io/gorm"
)

type ListData struct {
	db *gorm.DB
}

func NewListData(db *gorm.DB) *ListData {
	return &ListData{db: db}
}

func (d *ListData) FindAllVerUnique(list *models.CatalogList) ([]models.ListDataUnique, error) {
	var listDataUnique []models.ListDataUnique
	err := d.db.Table(list.TableName).Find(&listDataUnique).Error
	return listDataUnique, err
}
