package space

import (
	"fmt"
	"termonopoly/game"
)

type Corner struct {
	BaseSpace
}

func (c *Corner) Print() {
	fmt.Printf("%s\n", c.Name)
}

func CornerFactory(row []string) game.SpaceInterface {
	// fmt.Printf("Corner Item\n")
	switch row[1] {
	case "Go":
		g := &Go{}
		g.SetName(row[1])
		return g
	case "Go To Jail":
		j := &GoToJail{}
		j.SetName(row[1])
		return j
	case "Free Parking":
		p := &GoToJail{}
		p.SetName(row[1])
		return p
	case "Jail":
		j := &Jail{}
		j.SetName(row[1])
		return j

	}
	return &Corner{}
}

type Jail struct {
	Corner
}

type GoToJail struct {
	Corner
}

type Parking struct {
	Corner
}

type Go struct {
	Corner
}

func (Go) OnPass(pl game.OwnerInterface) {
	pl.PassGo()
}
