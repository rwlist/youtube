package proto

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
	Info(listID string) (ListInfo, error)

	//gjrpc:method items
	Items(listID string) (ListItems, error)

	//gjrpc:method pageItems
	PageItems(req PageRequest) (ListItems, error)

	//gjrpc:method executeQuery
	ExecuteQuery(query Query) (QueryResponse, error)
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

type ListItem struct {
	YoutubeID string
	Title     string
	Author    string
	ChannelID string
	ItemID    uint
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
