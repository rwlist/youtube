package lists

import "github.com/rwlist/youtube/internal/models"

type debugger interface {
	debugInfo() debugInfo
}

type debugInfo struct {
	// don't want to have *gorm.DB inside, so spew will be working ok
	storage    func() *Storage
	catalog    *models.CatalogList
	commonOpts *CommonOpts
}
