package ytsync

import "google.golang.org/api/youtube/v3"

type CursorOpts struct {
	// UserID for fetch the playlist for
	UserID uint

	// Return pages of this size
	PageSize int64
}

type CursorCreator interface {
	CreateCursor(opts *CursorOpts) (Cursor, error)
}

type Cursor interface {
	Next() (*youtube.PlaylistItemListResponse, error)
}
