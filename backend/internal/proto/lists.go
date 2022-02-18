package proto

//gjrpc:service catalog
type ListsCatalog interface {
	//gjrpc:method all
	All() (*AllLists, error)
}

type AllLists struct {
	Lists []ListInfo
}
