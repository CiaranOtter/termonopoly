package space

import (
	"log"
	"strconv"
	"termonopoly/game"
)

type BaseSpace struct {
	game.SpaceInterface
	Name  string
	Rent  int
	Price int
	Group string

	prev game.SpaceInterface
	next game.SpaceInterface
}

func SpaceFactory(row []string) game.SpaceInterface {
	// fmt.Println(row)

	var space game.SpaceInterface
	switch row[2] {
	case "Corners":
		space = CornerFactory(row)
	case "Railroads":
		space = NewRailRoad(row)
	case "Utilities":
		space = NewUtilities(row)
	case "Cards":
		space = CardSpaceFactory(row)
	case "Tax":
		space = NewTax(row)
	default:
		space = NewProperty(row)
	}

	space.SetGroup(row[2])
	return space
}

func (b *BaseSpace) SetGroup(group string) {
	b.Group = group
}

func (b *BaseSpace) GetGroup() string {
	return b.Group
}

func (b *BaseSpace) GetNext() game.SpaceInterface {
	return b.next
}

func (b *BaseSpace) GetPrev() game.SpaceInterface {
	return b.prev
}

func (b *BaseSpace) SetNext(p game.SpaceInterface) {

	b.next = p

}

func (b *BaseSpace) SetPrev(p game.SpaceInterface) {

	b.prev = p
}

func (b *BaseSpace) SetName(name string) {
	b.Name = name
}

func (b *BaseSpace) SetPrice(price string) {
	a, err := strconv.Atoi(price)

	if err != nil {
		log.Fatal(err)
	}

	b.Price = a
}

func (b *BaseSpace) SetRent(rent string) {
	a, err := strconv.Atoi(rent)

	if err != nil {
		log.Fatal(err)
	}

	b.Rent = a
}

func (b *BaseSpace) Print() {

}

func (b *BaseSpace) OnLand(pl game.OwnerInterface) {

}

func (b *BaseSpace) OnPass(pl game.OwnerInterface) {

}
