package termonopoly

import (
	"fmt"
	"math/rand"
)

type Node interface {
	Next() Node
	SetNext(Node)
	GetName() string
	GetType() string
	OnPass(player *Player)
	OnLand(player *Player)
	Print()
}

type Space struct {
	Index int
	Name  string
	Type  int
	Next  Node
}

var Start Node

var Properties map[string][]*Property
var ChanceCards []*ChanceCard
var CommunityCards []*CommunityCard
var UsedChanceCards []*ChanceCard
var UsedCommunityCards []*CommunityCard
var Players []*Player

func InitPlayers(count int) []*Player {
	var players = make([]*Player, count)

	for i := 0; i < count; i++ {
		player := InitPlayer(fmt.Sprintf("Player %d", i))
		players[i] = &player
	}

	return players
}

func ListAvailable() {
	for colour, group := range Properties {
		fmt.Printf("%s: \n", colour)
		for _, prop := range group {
			if prop.Owner == nil {
				prop.Print()
			}
		}
	}
}

func DrawChanceCard() *ChanceCard {
	i := rand.Intn(len(ChanceCards))

	card := ChanceCards[i]

	ChanceCards = append(ChanceCards[:i], ChanceCards[i+1:]...)
	UsedChanceCards = append(UsedChanceCards, card)

	if len(ChanceCards) == 0 {
		ChanceCards = UsedChanceCards
		UsedChanceCards = make([]*ChanceCard, 0)
	}

	return card
}

func DrawCommunityCard() *CommunityCard {
	i := rand.Intn(len(CommunityCards))

	card := CommunityCards[i]

	CommunityCards = append(CommunityCards[:i], CommunityCards[i+1:]...)
	UsedCommunityCards = append(UsedCommunityCards, card)

	if len(CommunityCards) == 0 {
		CommunityCards = UsedCommunityCards
		UsedCommunityCards = make([]*CommunityCard, 0)
	}

	return card
}

func RollDice() (int, bool) {
	dice1 := rand.Intn(5) + 1
	dice2 := rand.Intn(5) + 1

	total := dice1 + dice2

	fmt.Printf("====== Dice Roll =======\n")
	if dice1 == dice2 {
		fmt.Printf("Double rolled! -> %d\n", dice1)
	}
	fmt.Printf("Roll total: %d\n", total)
	fmt.Printf("========================\n")

	return total, dice1 == dice2
}
