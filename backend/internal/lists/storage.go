package lists

import (
	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/proto"
	"gorm.io/gorm"
)

// TODO: can storage be a generic interface, supporting different storage backends?
type Storage struct {
	db      *gorm.DB
	catalog *models.CatalogList
}

func NewStorage(db *gorm.DB, catalog *models.CatalogList) *Storage {
	return &Storage{
		db:      db,
		catalog: catalog,
	}
}

func (s Storage) Transaction(f func(s *Storage) error) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		s.db = tx
		return f(&s)
	})
}

func (s *Storage) Catalog() *models.CatalogList {
	return s.catalog
}

func (s *Storage) CreateTable(model interface{}) error {
	return s.db.Table(s.catalog.TableName).Migrator().CreateTable(model)
}

func (s Storage) Automigrate(model interface{}) error {
	return s.db.Table(s.catalog.TableName).Migrator().AutoMigrate(model)
}

func (s *Storage) FindAll(destArr interface{}) error {
	return s.db.Table(s.catalog.TableName).Order("xord ASC").Find(destArr).Error
}

func (s *Storage) FindByPageRequest(req proto.PageRequest, destArr interface{}) error {
	// TODO: how to make it work faster?
	return s.db.
		Table(s.catalog.TableName).
		Order("xord ASC").
		Offset(req.Offset).
		Limit(req.Limit).
		Find(destArr).
		Error
}

func (s *Storage) FindByPrimaryKey(destArr interface{}, primaryKeys interface{}) error {
	return s.db.Table(s.catalog.TableName).Find(destArr, primaryKeys).Error
}

func (s *Storage) FindByYoutubeID(destArr interface{}, ids []string) error {
	return s.db.Table(s.catalog.TableName).Where("youtube_id IN ?", ids).Find(destArr).Error
}

func (s *Storage) FirstByYoutubeID(dest interface{}, id string) error {
	return s.db.Table(s.catalog.TableName).Where("youtube_id = ?", id).Find(dest).Error
}

func (s *Storage) OrderLimit(res interface{}, xord string, cnt int) error {
	return s.db.Table(s.catalog.TableName).Order("xord ASC").Limit(cnt).Find(res, "xord >= ?", xord).Error
}

func (s *Storage) UpdateXord(id uint, xord string) error {
	return s.db.Table(s.catalog.TableName).Where("item_id = ?", id).UpdateColumn("xord", xord).Error
}

func (s *Storage) Insert(model interface{}) error {
	return s.db.Table(s.catalog.TableName).Create(model).Error
}

func (s *Storage) Save(model interface{}) error {
	return s.db.Table(s.catalog.TableName).Save(model).Error
}

func (s *Storage) CountAll() (int, error) {
	var count int64
	err := s.db.Table(s.catalog.TableName).Count(&count).Error
	return int(count), err
}
