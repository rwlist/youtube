package proto

import "google.golang.org/api/youtube/v3"

//gjrpc:service youtube
type Youtube interface {
	//gjrpc:method playlists
	Playlists() (*Playlists, error)

	//gjrpc:method liked
	Liked() (*PlaylistItems, error)
}

type Playlists struct {
	Response *youtube.PlaylistListResponse
}

type PlaylistItems struct {
	Response *youtube.PlaylistItemListResponse
}
