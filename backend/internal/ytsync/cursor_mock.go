package ytsync

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/rwlist/youtube/internal/models"
	"google.golang.org/api/youtube/v3"
)

type MockCursor struct {
	Items []models.YoutubeData
	Rand  *rand.Rand

	opts *CursorOpts
}

func (l MockCursor) CreateCursor(opts *CursorOpts) (Cursor, error) {
	l.opts = opts
	return &l, nil
}

func (l *MockCursor) GenItem(i int) models.YoutubeData {
	now := time.Unix(int64(1892160000+i*60*60), 0)
	return models.YoutubeData{
		YoutubeID:   "youtube" + strconv.Itoa(i),
		Title:       "Video #" + strconv.Itoa(l.Rand.Int()),
		Author:      "Author is the same",
		ChannelID:   "Discovery channel",
		PublishedAt: &now,
	}
}

func (l *MockCursor) Shuffle() {
	l.Rand.Shuffle(len(l.Items), func(i, j int) {
		l.Items[i], l.Items[j] = l.Items[j], l.Items[i]
	})
}

func (l *MockCursor) Next() (*youtube.PlaylistItemListResponse, error) {
	var cur []models.YoutubeData
	if len(l.Items) <= int(l.opts.PageSize) {
		cur = l.Items
		l.Items = nil
	} else {
		cur = l.Items[:l.opts.PageSize]
		l.Items = l.Items[l.opts.PageSize:]
	}

	resp := &youtube.PlaylistItemListResponse{
		PageInfo: &youtube.PageInfo{
			ResultsPerPage: l.opts.PageSize,
		},
	}
	for _, item := range cur {
		item := item
		resp.Items = append(resp.Items, models.ToYoutubeItem(&item))
	}

	return resp, nil
}
