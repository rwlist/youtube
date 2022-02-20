package lists

import (
	"fmt"

	"github.com/rwlist/youtube/internal/global"
	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/proto"
	"github.com/rwlist/youtube/internal/repos"
	"gorm.io/gorm"
)

const (
	listLiked = "liked"
)

type Catalog struct {
	repo      *repos.CatalogLists
	defaultDB *gorm.DB
	dir       *global.Directory
}

func NewCatalog(repo *repos.CatalogLists, defaultDB *gorm.DB, dir *global.Directory) *Catalog {
	return &Catalog{
		repo:      repo,
		defaultDB: defaultDB,
		dir:       dir,
	}
}

// UserHook creates some common lists for user on login
func (c *Catalog) UserHook(user *models.User) error {
	// check liked list
	if ok, err := c.listExists(user.ID, listLiked); err != nil {
		return err
	} else if ok {
		return nil
	}

	// create liked list
	entry := &models.CatalogList{
		UserID:   user.ID,
		ListID:   listLiked,
		ListName: "Liked videos",
		ListType: proto.ListTypeExternal,
	}

	return c.repo.Create(entry, func(db *gorm.DB) error {
		engine, err := c.buildEngine(db, entry)
		if err != nil {
			return err
		}

		return engine.InitStorage()
	})
}

func (c *Catalog) listExists(userID uint, listID string) (bool, error) {
	_, err := c.repo.GetByListID(userID, listID)
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *Catalog) UserLists(userID uint) ([]Engine, error) {
	res, err := c.repo.FindUserLists(userID)
	if err != nil {
		return nil, err
	}

	lists := make([]Engine, len(res))
	for i := range res {
		lists[i], err = c.buildEngine(c.defaultDB, &res[i])
		if err != nil {
			return nil, err
		}
	}
	return lists, nil
}

func (c *Catalog) GetList(userID uint, listID string) (Engine, error) {
	catalog, err := c.repo.GetByListID(userID, listID)
	if err != nil {
		return nil, err
	}

	return c.buildEngine(c.defaultDB, catalog)
}

func (c *Catalog) buildEngine(db *gorm.DB, list *models.CatalogList) (Engine, error) {
	storage := NewStorage(db, list, c)
	if list.ListType == proto.ListTypeExternal && list.ListID == listLiked {
		return NewLikedEngine(storage, c.dir), nil
	}

	return nil, fmt.Errorf("list %s not supported, list type: %s", list.ListID, list.ListType)
}

func (c *Catalog) engineUpdateCatalog(tx *gorm.DB, ct *models.CatalogList) error {
	return c.repo.Update(tx, ct)
}
