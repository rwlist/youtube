package global

type LikedSync interface {
	Sync(id string, engine LikedEngine) error
}

type Directory struct {
	LikedSync LikedSync
}
