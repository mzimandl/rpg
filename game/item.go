package game

type Item struct {
	Entity
}

func NewSword(p Pos) *Item {
	item := &Item{Entity{p, 's', "Sword"}}
	return item
}

func NewHelmet(p Pos) *Item {
	item := &Item{Entity{p, 'h', "Helmet"}}
	return item
}
