package models

import (
	"github.com/rwlist/youtube/internal/proto"
	"google.golang.org/api/youtube/v3"
	"time"
)

// Model is generic header for all objects in the list. Should be embedded in user structs.
type Model struct {
	// permanent id
	ItemID    uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// actually can be non unique, but still changes frequently
	Xord string `gorm:"not null;unique"`
}

func XordModel(xord string) Model {
	return Model{
		Xord: xord,
	}
}

func (m Model) Meta() proto.Meta {
	return proto.Meta{
		ItemID:    m.ItemID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Xord:      m.Xord,
	}
}

// LikedModel is used for automigration.
type LikedModel struct {
	Model
	LikedData
}

type LikedData struct {
	// generic youtube info, can only update in rare cases
	YoutubeID string `gorm:"not null;unique"`
	Title     string `gorm:"not null"`
	Author    string `gorm:"not null"`
	ChannelID string `gorm:"not null"`

	// additional info
	PublishedAt *time.Time // when added to original playlist
}

func LikedDataFromYoutube(data *YoutubeData) *LikedData {
	res := LikedData(*data)
	return &res
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

func FromYoutubeItem(item *youtube.PlaylistItem) (*YoutubeData, error) {
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
