package proto

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
	ID   string
	Name string

	// external list: cannot be modified, only can be synced
	// custom list: can be modified, cannot be synced
	// virtual list: fully handled by the client
	// type ListType = 'external' | 'custom' | 'virtual'
	ListType string
}

type ListItems struct {
	Items []ListItem
}

type ListItem struct {
	YoutubeID string
	Title     string
	Author    string
	ChannelID string
	ItemID    string
	Xord      string
}
