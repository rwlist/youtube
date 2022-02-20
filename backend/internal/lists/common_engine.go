package lists

import (
	"github.com/rwlist/youtube/internal/global"
	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/proto"
)

// with unique object_id, no joins

type CommonEngine struct {
	storage *Storage
	dir     *global.Directory
	opts    *CommonOpts
}

type CommonOpts struct {
	newModel     func() interface{}
	newArray     func() interface{}
	convertArray func(interface{}) ([]proto.ListItem, error)
	initMeta     models.CatalogMeta
}

func NewCommonEngine(storage *Storage, dir *global.Directory, opts *CommonOpts) CommonEngine {
	return CommonEngine{
		storage: storage,
		dir:     dir,
		opts:    opts,
	}
}

func (e *CommonEngine) InitStorage() error {
	ct := *e.storage.Catalog()
	meta := e.opts.initMeta
	ct.Meta = &meta
	err := e.storage.UpdateCatalog(ct)
	if err != nil {
		return err
	}

	return e.storage.CreateTable(e.opts.newModel())
}

func (e *CommonEngine) AutoMigrate() error {
	return e.storage.AutoMigrate(e.opts.newModel())
}

func (e *CommonEngine) CatalogEntry() (*models.CatalogList, error) {
	return e.storage.Catalog(), nil
}

func (e *CommonEngine) Info() (*proto.ListInfo, error) {
	info := e.storage.Catalog().ToInfo()
	itemsCount, err := e.storage.CountAll()
	if err != nil {
		return nil, err
	}
	info.ItemsCount = itemsCount
	return info, nil
}

func (e *CommonEngine) ListItems() ([]proto.ListItem, error) {
	items := e.opts.newArray()
	err := e.storage.FindAll(items)
	if err != nil {
		return nil, err
	}

	return e.opts.convertArray(items)
}

func (e *CommonEngine) PageItems(req *proto.PageRequest) ([]proto.ListItem, error) {
	items := e.opts.newArray()
	err := e.storage.FindByPageRequest(req, items)
	if err != nil {
		return nil, err
	}

	return e.opts.convertArray(items)
}

// Pass "" to insert in the beginning.
func (e *CommonEngine) xordForInsert(afterXord string) (string, error) {
	cnt := 2
	if afterXord != "" {
		cnt = 3
	}

	var start []models.Model
	err := e.storage.OrderLimit(&start, afterXord, cnt)
	if err != nil {
		return "", err
	}

	// <afterXord> <new> <[0]> <[1]>
	newXord := ""
	if len(start) == 0 {
		newXord = splitXord("", "")
	} else if afterXord == start[0].Xord && len(start) == 1 {
		newXord = splitXord(start[0].Xord, "")
	} else {
		if afterXord == start[0].Xord {
			start = start[1:]
		}
		newXord = start[0].Xord
		err = e.shiftXords(start)
		if err != nil {
			return "", err
		}
	}

	return newXord, nil
}

// Shift first element to the right.
// If <xord1> <xord2> is passed, then first element will be shifted to split(xord1, xord2)
func (e *CommonEngine) shiftXords(nxt2 []models.Model) error {
	if len(nxt2) == 0 {
		return nil
	}
	l := nxt2[0].Xord
	r := ""
	if len(nxt2) == 2 {
		r = nxt2[1].Xord
	}

	newXord := splitXord(l, r)
	return e.storage.UpdateXord(nxt2[0].ItemID, newXord)
}

func (e *CommonEngine) MoveAfter(xord string, itemID uint) (newXord string, err error) {
	newXord, err = e.xordForInsert(xord)
	if err != nil {
		return "", err
	}

	return newXord, e.storage.UpdateXord(itemID, newXord)
}

func (e *CommonEngine) debugInfo() debugInfo {
	return debugInfo{
		storage:    func() *Storage { return e.storage },
		catalog:    e.storage.Catalog(),
		commonOpts: e.opts,
	}
}
