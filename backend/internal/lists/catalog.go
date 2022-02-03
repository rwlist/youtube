package lists

import (
	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/proto"
	"github.com/rwlist/youtube/internal/repos"
	"gorm.io/gorm"
)

const (
	listLiked = "liked"
)

type Catalog struct {
	repo     *repos.CatalogLists
	listData *repos.ListData
}

func NewCatalog(repo *repos.CatalogLists, listData *repos.ListData) *Catalog {
	return &Catalog{
		repo:     repo,
		listData: listData,
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
		return c.createDataTable(db, entry)
	})
}

func (c *Catalog) createDataTable(db *gorm.DB, list *models.CatalogList) error {
	return db.Table(list.TableName).Migrator().CreateTable(&models.ListDataUnique{})
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

func (c *Catalog) UserLists(userID uint) ([]proto.ListInfo, error) {
	res, err := c.repo.FindUserLists(userID)
	if err != nil {
		return nil, err
	}

	lists := make([]proto.ListInfo, len(res))
	for i, l := range res {
		lists[i] = proto.ListInfo{
			ID:       l.ListID,
			Name:     l.ListName,
			ListType: l.ListType,
		}
	}
	return lists, nil
}

func (c *Catalog) ViewList(userID uint, listID string) ([]proto.ListItem, error) {
	list, err := c.repo.GetByListID(userID, listID)
	if err != nil {
		return nil, err
	}

	items, err := c.listData.FindAllVerUnique(list)
	if err != nil {
		return nil, err
	}

	res := make([]proto.ListItem, len(items))
	for i, item := range items {
		res[i] = proto.ListItem{
			YoutubeID: item.YoutubeID,
			Title:     item.Title,
			Author:    item.Author,
			ChannelID: item.ChannelID,
			ItemID:    item.ItemID,
			Xord:      item.Xord,
		}
	}
	return res, nil
}
