package game

type Game struct {
	Jail SpaceInterface
}

func NewGame() *Game {
	return &Game{}
}
