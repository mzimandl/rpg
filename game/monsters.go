package game

type Monster struct {
	Character
}

func NewRat(pos Pos) *Monster {
	monster := &Monster{}
	monster.Pos = pos
	monster.Rune = 'R'
	monster.Name = "Rat"
	monster.Hitpoints = 5
	monster.Strength = 5
	monster.Speed = 2.0
	monster.ActionPoints = 0.0
	return monster
}

func NewSpider(pos Pos) *Monster {
	monster := &Monster{}
	monster.Pos = pos
	monster.Rune = 'S'
	monster.Name = "Spider"
	monster.Hitpoints = 10
	monster.Strength = 10
	monster.Speed = 1.0
	monster.ActionPoints = 0.0
	return monster
}

func (m *Monster) Update(level *Level) {
	m.ActionPoints += m.Speed
	positions := level.astar(m.Pos, level.Player.Pos)
	for m.ActionPoints >= 1 {
		if positions != nil {
			m.ActionPoints--
			next := positions[0]
			if next == level.Player.Pos {
				event := m.Attack(&level.Player.Character)
				level.addEvent(event)
				if !level.Player.IsAlive() {
					level.addEvent("DED")
				}
			} else {
				moved := m.Move(level, next)
				if moved {
					positions = positions[1:]
				} else {
					m.Pass()
				}
			}
		} else {
			m.Pass()
		}
	}
}

func (m *Monster) Pass() {
	m.ActionPoints = 0
}

func (m *Monster) Move(level *Level, next Pos) bool {
	_, exists := level.AliveMonstersPos[next]
	if exists {
		return false
	}

	delete(level.AliveMonstersPos, m.Pos)
	m.Pos = next
	level.AliveMonstersPos[next] = m
	return true
}

func (m *Monster) Die(level *Level) {
	level.Monsters = append(level.Monsters, m)
	delete(level.AliveMonstersPos, m.Pos)
}
