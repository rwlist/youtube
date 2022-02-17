package models

import (
	"crypto/rand"
	"encoding/json"
	"github.com/rwlist/youtube/internal/proto"
	"gorm.io/gorm"
	"time"
)

type CatalogMeta struct {
	// empty -> no object id
	ObjectIDField    string
	IsUniqueObjectID bool
}

type CatalogList struct {
	TableName string `gorm:"type:char(15);primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID uint   `gorm:"type:bigint REFERENCES users(id);not null;index:idx_list_id,unique"`
	ListID string `gorm:"index:idx_list_id,unique;not null"`

	ListName string
	ListType proto.ListType `gorm:"not null"`

	MetaJSON []byte       `gorm:"type:jsonb"`
	Meta     *CatalogMeta `gorm:"-"`
}

func (l *CatalogList) GenerateTableName() error {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	const length = 10

	l.TableName = "list_"
	// use secure random to generate a random string
	arr := make([]byte, length)
	_, err := rand.Read(arr)
	if err != nil {
		return err
	}

	for i := 0; i < length; i++ {
		l.TableName += string(alphabet[arr[i]%byte(len(alphabet))])
	}
	return nil
}

func (l *CatalogList) ToInfo() *proto.ListInfo {
	return &proto.ListInfo{
		ID:       l.ListID,
		Name:     l.ListName,
		ListType: l.ListType,
	}
}

func (l *CatalogList) AfterLoad() error {
	l.Meta = &CatalogMeta{}
	if l.MetaJSON == nil {
		return nil
	}

	if err := json.Unmarshal(l.MetaJSON, l.Meta); err != nil {
		return err
	}
	return nil
}

func (l *CatalogList) BeforeSave() error {
	if l.Meta == nil {
		l.Meta = &CatalogMeta{}
	}
	var err error
	l.MetaJSON, err = json.Marshal(l.Meta)
	return err
}
