package game

type ItemType int

const (
	Weapon ItemType = iota
	Helmet
	Other
)

type Item struct {
	Entity
	Typ   ItemType
	Power float64
}

func NewSword(p Pos) *Item {
	item := &Item{Entity{p, 's', "Sword"}, Weapon, 2.0}
	return item
}

func NewHelmet(p Pos) *Item {
	item := &Item{Entity{p, 'h', "Helmet"}, Helmet, 0.1}
	return item
}
