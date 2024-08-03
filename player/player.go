package player

import (
	"fmt"
	"termonopoly/controls"
	"termonopoly/game"
	// "termonopoly/space"
)

type Player struct {
	game.OwnerInterface
	Name       string
	Money      int
	Space      game.SpaceInterface
	Properties map[string]([]game.SpaceInterface)
	Game       *game.Game
	InJail     bool
	JailCount  int
}

const (
	FORWARD  = 0
	BACKWARD = 1
)

func (p *Player) AddProperty(prop game.SpaceInterface) {

	group, exs := p.Properties[prop.GetGroup()]

	if !exs {
		p.Properties[prop.GetGroup()] = make([]game.SpaceInterface, 1)
		p.Properties[prop.GetGroup()][0] = prop
	} else {
		p.Properties[prop.GetGroup()] = append(group, prop)
	}

}

func (p *Player) Roll() {

}

func (p *Player) Move(count int, dir int) {
	for i := 0; i < count; i++ {
		if dir == FORWARD {
			p.Space = p.Space.GetNext()
		} else {
			p.Space.GetPrev()
		}

		p.Space.OnPass(p)
	}

	fmt.Printf("Landed on:\n")
	p.Space.Print()

	p.Space.OnLand(p)
	// p.OnLand(p)
}

func (p *Player) PassGo() {
	p.Money += 200
}

func (p *Player) GoToJail() {
	p.Space = p.Game.Jail
	p.InJail = true
}

func (p *Player) OfferProperty(prop game.SpaceInterface) {
	cost, canAfford := prop.Afford(p.Money)

	if canAfford {
		question := fmt.Sprintf("This property costs $%d. You have $%d to your name.\n would you like to buy it?", cost, p.Money)

		if controls.YesNo(question) {
			p.Money -= cost
			p.AddProperty(prop)
			prop.SetOwner(p)
		}
	}

	fmt.Printf("You can not afford this property: $%d\n", cost)
}

func NewPlayer(name string, start game.SpaceInterface, startMoney int) *Player {
	return &Player{
		Name:       name,
		Space:      start,
		Money:      startMoney,
		Properties: make(map[string][]game.SpaceInterface),
	}
}
