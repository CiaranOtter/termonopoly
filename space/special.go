package space

import (
	"strings"
	"termonopoly/game"
)

type Tax struct {
	BaseSpace
}

func (t *Tax) Print() {
	// fmt.Printf("Tax space\n")
}

func NewTax(row []string) *Tax {
	// fmt.Printf("Tax space\n")
	return &Tax{}
}

type CardSpace struct {
	BaseSpace
}

func CardSpaceFactory(row []string) game.SpaceInterface {
	// fmt.Printf("Card space\n")

	if strings.Contains(row[1], "Chance") {
		return NewChance(row)
	}

	if strings.Contains(row[1], "Community Chest") {
		return NewCommunityChest(row)
	}

	return &BaseSpace{}
}

type CommunityChest struct {
	CardSpace
}

func (cc *CommunityChest) Print() {
	// fmt.Printf("Community chest space\n")
}

func NewCommunityChest(row []string) *CommunityChest {
	// fmt.Printf("Community chest item\n")
	return &CommunityChest{}
}

type ChanceCard struct {
	CardSpace
}

func (cc *ChanceCard) Print() {
	// fmt.Printf("Chance card space\n")
}

func NewChance(row []string) *ChanceCard {
	// fmt.Printf("Chance card item\n")
	return &ChanceCard{}
}

type RailRoads struct {
	BaseSpace
}

func (r *RailRoads) Print() {
	// fmt.Printf("Railroad: %s\n", r.Name)
}

func NewRailRoad(row []string) *RailRoads {
	// fmt.Printf("Rail road item\n")
	r := &RailRoads{}
	r.SetName(row[1])
	r.SetPrice(row[3])
	r.SetRent(row[6])

	return r
}

type UtilitiesSpace struct {
	BaseSpace
}

func (u *UtilitiesSpace) Print() {
	// fmt.Printf("Utilities: %s\n", u.Name)
}

func NewUtilities(row []string) *UtilitiesSpace {
	// fmt.Printf("Utilities space\n")
	u := &UtilitiesSpace{}
	u.SetName(row[1])
	u.SetPrice(row[3])
	u.SetRent(row[5])

	return u
}
