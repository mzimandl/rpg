package game

type Repository struct {
	Entity
	Items []*Item
}

type Storage struct {
	Repository
	Locked bool
}

func NewChest(pos Pos) *Storage {
	chest := &Storage{}
	chest.Name = "Chest"
	chest.Rune = '='
	chest.Pos = pos
	chest.Items = append(chest.Items, NewSword(pos))
	chest.Locked = true
	return chest
}
