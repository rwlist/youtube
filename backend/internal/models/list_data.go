package models

import (
	"time"
)

type ListDataUnique struct {
	// permanent id
	ItemID    uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// actually can be non unique, but still changes frequently
	Xord string `gorm:"not null;unique"`

	// generic youtube info, can only update in rare cases
	YoutubeID string `gorm:"not null;unique"`
	Title     string `gorm:"not null"`
	Author    string `gorm:"not null"`
	ChannelID string `gorm:"not null"`
}
