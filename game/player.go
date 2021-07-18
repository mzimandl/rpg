package game

type Player struct {
	Character
}

func NewPlayer(pos Pos) *Player {
	player := &Player{}
	player.Pos = pos
	player.Rune = '@'
	player.Name = "Player"
	player.Hitpoints = 20
	player.Strength = 20
	player.Speed = 1.0
	player.ActionPoints = 0.0
	player.SightRange = 10
	return player
}

func (player *Player) Move(to Pos, level *Level) {
	player.Pos = to
}
