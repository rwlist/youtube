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

//gjrpc:service lists
type ListService interface {
	//gjrpc:method all
	All() (AllLists, error)

	//gjrpc:method listInfo
	ListInfo(listID string) (ListInfo, error)

	//gjrpc:method listItems
	ListItems(listID string) (ListItems, error)
}

type AllLists struct {
	Lists []ListInfo
}

type ListInfo struct {
	// Short and unique for user, e.g. "liked"
	ID string

	// Human-readable name of the list
	Name string

	// Type of the list
	ListType ListType
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
