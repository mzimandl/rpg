package game

type Repository struct {
	Entity
	Items []*Item
}

type Storage struct {
	Repository
	Locked bool
}

type StorageConf struct {
	items  []*Item
	locked bool
}

func NewChest(pos Pos, conf *StorageConf) *Storage {
	chest := &Storage{}
	chest.Name = "Chest"
	chest.Rune = '='
	chest.Pos = pos
	chest.Items = conf.items
	chest.Locked = conf.locked
	return chest
}
