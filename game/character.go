package game

import "strconv"

type Character struct {
	Entity
	Hitpoints    int
	Strength     int
	Speed        float64
	ActionPoints float64
	SightRange   int
	Items        []*Item
}

func (ch *Character) IsAlive() bool {
	return ch.Hitpoints > 0
}

func (ch *Character) Attack(ch2 *Character) string {
	ch2.Hitpoints -= ch.Strength

	if ch2.IsAlive() {
		return ch.Name + " hits " + ch2.Name + " causing damage " + strconv.Itoa(ch.Strength)
	} else {
		return ch.Name + " killed " + ch2.Name
	}
}
