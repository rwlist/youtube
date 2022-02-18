package proto

import "time"

type ListType string

const (
	// ListTypeExternal cannot be modified, only can be synced
	ListTypeExternal ListType = "external"
	// ListTypeCustom can be modified, cannot be synced
	ListTypeCustom ListType = "custom"
	// ListTypeVirtual fully handled by the client
	ListTypeVirtual ListType = "virtual"
)

//gjrpc:service list
type ListService interface {
	//gjrpc:method info
	Info(listID string) (*ListInfo, error)

	//gjrpc:method items
	Items(listID string) (*ListItems, error)

	//gjrpc:method pageItems
	PageItems(req *PageRequest) (*ListItems, error)

	//gjrpc:method executeQuery
	ExecuteQuery(query *Query) (*QueryResponse, error)
}

type ListInfo struct {
	// Short and unique for user, e.g. "liked"
	ID string

	// Human-readable name of the list
	Name string

	// Type of the list
	ListType ListType

	// Number of items in the list
	ItemsCount int
}

type ListItems struct {
	Items []ListItem
}

// ListItem is any struct that have Meta embedded in it, must be present in serialized JSON too.
type ListItem interface {
	Meta() Meta
}

// ItemLiked is a type declaration for usage in typescript.
type ItemLiked struct {
	ItemID uint
	Xord   string

	YoutubeID string
	Title     string
	Author    string
	ChannelID string
}

type Meta struct {
	ItemID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Xord      string
}

type Query struct {
	ListID string
	Query  string
}

type QueryResponse struct {
	Status string
	Object interface{} `json:",omitempty"`
}

type PageRequest struct {
	ListID string

	// TODO: add query by xord
	Offset int
	Limit  int
}
