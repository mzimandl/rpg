package game

import "strconv"

type Character struct {
	Repository

	Hitpoints    int
	Strength     int
	Speed        float64
	ActionPoints float64
	SightRange   int

	Helmet *Item
	Weapon *Item
	Armor  *Item
}

func (c *Character) IsAlive() bool {
	return c.Hitpoints > 0
}

func (c *Character) Attack(cToAttack *Character) string {
	attackPower := c.Strength
	if c.Weapon != nil {
		attackPower = int(float64(attackPower) * c.Weapon.Power)
	}
	damage := attackPower
	if cToAttack.Helmet != nil {
		damage = int(float64(damage) * (1.0 - cToAttack.Helmet.Power))
	}
	if cToAttack.Armor != nil {
		damage = int(float64(damage) * (1.0 - cToAttack.Armor.Power))
	}

	cToAttack.Hitpoints -= damage
	if cToAttack.IsAlive() {
		return c.Name + " hits " + cToAttack.Name + " causing damage " + strconv.Itoa(damage)
	} else {
		return c.Name + " killed " + cToAttack.Name + " causing damage " + strconv.Itoa(damage)
	}
}

func (c *Character) TakeItem(level *Level, itemToMove *Item) bool {
	items := level.Items[c.Pos]
	for i, item := range items {
		if item == itemToMove {
			level.Items[c.Pos] = append(items[:i], items[i+1:]...)
			c.Items = append(c.Items, itemToMove)
			return true
		}
	}
	return false
}

func (c *Character) DropItem(level *Level, itemToMove *Item) bool {
	for i, item := range c.Items {
		if item == itemToMove {
			item.Pos = c.Pos
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			level.Items[c.Pos] = append(level.Items[c.Pos], item)
			return true
		}
	}
	return false
}

func (c *Character) StoreItem(level *Level, itemToMove *Item) bool {
	storage := level.Storages[c.Pos]
	if storage != nil && !storage.Locked {
		for i, item := range c.Items {
			if item == itemToMove {
				item.Pos = c.Pos
				c.Items = append(c.Items[:i], c.Items[i+1:]...)
				storage.Items = append(storage.Items, item)
				return true
			}
		}
	}
	return false
}

func (c *Character) WithdrawItem(level *Level, itemToMove *Item) bool {
	storage := level.Storages[c.Pos]
	if storage != nil && !storage.Locked {
		for i, item := range storage.Items {
			if item == itemToMove {
				storage.Items = append(storage.Items[:i], storage.Items[i+1:]...)
				c.Items = append(c.Items, itemToMove)
				return true
			}
		}
	}
	return false
}

func (c *Character) Equip(itemToEquip *Item) bool {
	for i, item := range c.Items {
		if item == itemToEquip {
			var replace *Item

			switch itemToEquip.Typ {
			case Helmet:
				if c.Helmet != nil {
					replace = c.Helmet
				}
				c.Helmet = itemToEquip
			case Weapon:
				if c.Weapon != nil {
					replace = c.Weapon
				}
				c.Weapon = itemToEquip
			case Armor:
				if c.Armor != nil {
					replace = c.Armor
				}
				c.Armor = itemToEquip
			default:
				return false
			}

			if replace != nil {
				c.Items[i] = replace
			} else {
				c.Items = append(c.Items[:i], c.Items[i+1:]...)
			}
			return true
		}
	}
	return false
}

func (c *Character) Strip(itemToStrip *Item) bool {
	switch itemToStrip {
	case c.Helmet:
		c.Helmet = nil
	case c.Weapon:
		c.Weapon = nil
	case c.Armor:
		c.Armor = nil
	default:
		return false
	}

	c.Items = append(c.Items, itemToStrip)
	return true
}
