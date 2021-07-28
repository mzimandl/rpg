package game

type Storage struct {
	Entity
	Items []*Item
}

func NewChest(pos Pos) *Storage {
	chest := &Storage{}
	chest.Name = "Chest"
	chest.Rune = '='
	chest.Pos = pos
	chest.Items = append(chest.Items, NewSword(pos))
	return chest
}
