package lists

import (
	"fmt"

	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/proto"
	"gorm.io/gorm"
)

// TODO: can storage be a generic interface, supporting different storage backends?

type Storage struct {
	db             *gorm.DB
	txInProgress   bool
	catalog        *models.CatalogList
	catalogManager *Catalog
}

func NewStorage(db *gorm.DB, catalog *models.CatalogList, c *Catalog) *Storage {
	return &Storage{
		db:             db,
		catalog:        catalog,
		catalogManager: c,
	}
}

func (s Storage) Transaction(f func(s *Storage) error) error {
	if s.txInProgress {
		return fmt.Errorf("storage can't start embedded tx")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		s.db = tx
		s.txInProgress = true
		return f(&s)
	})
}

func (s *Storage) Catalog() *models.CatalogList {
	return s.catalog
}

func (s *Storage) CreateTable(model interface{}) error {
	return s.db.Table(s.catalog.TableName).Migrator().CreateTable(model)
}

func (s Storage) AutoMigrate(model interface{}) error {
	return s.db.Table(s.catalog.TableName).Migrator().AutoMigrate(model)
}

func (s *Storage) FindAll(destArr interface{}) error {
	return s.db.Table(s.catalog.TableName).Order("xord ASC").Find(destArr).Error
}

func (s *Storage) FindByPageRequest(req *proto.PageRequest, destArr interface{}) error {
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

func (s *Storage) FirstByKey(dest interface{}, key string, value interface{}) error {
	return s.db.Table(s.catalog.TableName).Where(key+" = ?", value).Find(dest).Error
}

func (s *Storage) FindByKeys(destArr interface{}, key string, values interface{}) error {
	return s.db.Table(s.catalog.TableName).Where(key+" IN ?", values).Find(destArr).Error
}

func (s *Storage) FirstByObjectID(dest interface{}, value interface{}) error {
	return s.FirstByKey(dest, s.catalog.Meta.ObjectIDField, value)
}

func (s *Storage) FindByObjectIDs(destArr interface{}, values interface{}) error {
	return s.FindByKeys(destArr, s.catalog.Meta.ObjectIDField, values)
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

func (s *Storage) UpdateCatalog(ct models.CatalogList) error {
	err := s.catalogManager.engineUpdateCatalog(s.db, &ct)
	if err != nil {
		return err
	}

	s.catalog = &ct
	return nil
}
