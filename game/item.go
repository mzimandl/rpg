package game

type Item struct {
	Entity
}

func NewSword(p Pos) *Item {
	item := &Item{}
	item.Pos = p
	item.Name = "Sword"
	item.Rune = 's'
	return item
}

func NewHelmet(p Pos) *Item {
	item := &Item{}
	item.Pos = p
	item.Name = "Helmet"
	item.Rune = 'h'
	return item
}
