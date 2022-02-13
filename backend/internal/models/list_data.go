package models

import (
	"google.golang.org/api/youtube/v3"
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

	// additional info
	PublishedAt *time.Time // when added to original playlist
}

func (u *ListDataUnique) UpdateData(data *YoutubeData) {
	u.YoutubeID = data.YoutubeID
	u.Title = data.Title
	u.Author = data.Author
	u.ChannelID = data.ChannelID
	u.PublishedAt = data.PublishedAt
}

type YoutubeData struct {
	// generic youtube info, can only update in rare cases
	YoutubeID string
	Title     string
	Author    string
	ChannelID string

	// additional info
	PublishedAt *time.Time // when added to original playlist
}

func ItemToData(item *youtube.PlaylistItem) (*YoutubeData, error) {
	publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	if err != nil {
		return nil, err
	}

	return &YoutubeData{
		YoutubeID:   item.ContentDetails.VideoId,
		Title:       item.Snippet.Title,
		Author:      item.Snippet.VideoOwnerChannelTitle,
		ChannelID:   item.Snippet.VideoOwnerChannelId,
		PublishedAt: &publishedAt,
	}, nil
}
