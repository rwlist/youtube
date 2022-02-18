package lists

import "github.com/rwlist/youtube/internal/proto"

type Engine interface {
	InitStorage() error
	Info() (*proto.ListInfo, error)
	ListItems() ([]proto.ListItem, error)
	ExecuteQuery(query string) (*proto.QueryResponse, error)
	PageItems(req *proto.PageRequest) ([]proto.ListItem, error)
}
