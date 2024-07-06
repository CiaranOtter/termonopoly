package termonopoly

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"termonopoly/termonopoly/comm"
)

type Player struct {
	Name        string
	Pos         Node
	Cash        int
	Properties  map[string]([]*Property)
	FullSets    []string
	InJail      bool
	DoubleCount int
	JailFree    int
	JailCount   int
	NetWorth    int
	Bankrupt    bool
	Stream      comm.Termonopoly_GameStreamClient
}

func (p *Player) Print() {
	fmt.Printf("======== Player ========\n")
	fmt.Printf("\n")
	fmt.Printf("Name:\t\t%s\n", p.Name)
	fmt.Printf("Cash:\t\t%d\n", p.Cash)
	fmt.Printf("Properties:\t%d\n", len(p.Properties))
	fmt.Printf("Net worth:\t%d\n", p.NetWorth)
	fmt.Printf("Full sets:\t%d\n", len(p.FullSets))

	if p.InJail {
		fmt.Printf("xxxxxxxxxxxx\n")
		fmt.Printf("%s is in jail\n", p.Name)
		fmt.Printf("xxxxxxxxxxxx\n")
	}

	fmt.Printf("Currently at:\n")
	fmt.Printf("\t%s\n", p.Pos.GetName())

	fmt.Printf("\n========================\n")
}

func (p *Player) SetBankrupt() {
	p.Bankrupt = true
}

func (p *Player) IsBankrupt() bool {
	return p.Bankrupt
}

func InitPlayer(name string) Player {
	return Player{
		Name:        name,
		Pos:         Start,
		Cash:        1500,
		InJail:      false,
		DoubleCount: 0,
		JailFree:    0,
		JailCount:   0,
		Properties:  make(map[string][]*Property),
		FullSets:    make([]string, 0),
		NetWorth:    1500,
		Bankrupt:    false,
	}
}

func (p *Player) SetInJail() {
	p.DoubleCount = 0
	p.JailCount = 0
	p.InJail = true
}

func (p *Player) SetFree() {
	p.DoubleCount = 0
	p.JailCount = 0
	p.InJail = false
}

func (p *Player) Spend(amount int) {
	p.Cash = p.Cash - amount
	p.NetWorth = p.NetWorth - amount

	fmt.Printf("Spent %d\n", amount)
	fmt.Printf("%d remaining\n", p.Cash)
	fmt.Print("\n")

	p.Stream.Send(&comm.Message{
		Type: comm.MessageType_ACTION,
		Data: &comm.Message_Act{
			&comm.Action{
				Type: comm.ActionType_SEND,
				Action: &comm.Action_Trans{
					&comm.Trans{
						Amount: int32(amount),
					},
				},
			},
		},
	})
}

func (p *Player) Receive(amount int) {
	p.Cash = p.Cash + amount

	p.NetWorth += amount

	fmt.Printf("Received %d\n", amount)
	fmt.Printf("%d in account\n", p.Cash)
	fmt.Printf("\n")

	p.Stream.Send(&comm.Message{
		Type: comm.MessageType_ACTION,
		Data: &comm.Message_Act{
			&comm.Action{
				Type: comm.ActionType_RECV,
				Action: &comm.Action_Trans{
					&comm.Trans{
						Amount: int32(amount),
					},
				},
			},
		},
	})
}

func (p *Player) AddProperty(prop *Property) {

	_, exists := p.Properties[prop.Family]

	if !exists {
		p.Properties[prop.Family] = make([]*Property, 0)
	}

	p.Properties[prop.Family] = append(p.Properties[prop.Family], prop)

	if len(p.Properties[prop.Family]) == p.Properties[prop.Family][0].TotalFamily {
		p.FullSets = append(p.FullSets, prop.Family)
	}

	p.NetWorth += prop.Price
}

func (p *Player) RemoveProperty(colour string, index int) {
	group, exists := p.Properties[colour]

	if !exists {
		return
	}

	if index < 0 || index >= len(group) {
		return
	}

	p.NetWorth -= group[index].Price

	group = append(group[:index], group[index+1:]...)

	if i := slices.Index(p.FullSets, colour); i != -1 {
		p.FullSets = append(p.FullSets[i:], p.FullSets[:i+1]...)
	}

}

func (p *Player) MoveUntil(family string, pass bool) {
	iter := p.Pos.Next()

	for strings.Compare(iter.GetType(), family) != 0 {

		if pass {
			iter.OnPass(p)
		}

		iter = iter.Next()
	}

	p.Pos = iter
	iter.OnLand(p)
}

func (p *Player) MoveDist(dist int, pass bool) {
	iter := p.Pos.Next()

	for i := 1; i < dist; i++ {
		if pass {
			iter.OnPass(p)
		}

		iter = iter.Next()
	}

	p.Pos = iter

	// move := &comm.Message{
	// 	Type: comm.MessageType_ACTION,
	// 	Data: &comm.Message_Act{
	// 		&comm.Action{
	// 			Action: &comm.Action_Move{
	// 				&comm.Move{
	// 					Dist: int32(dist),
	// 					Pass: pass,
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// p.Stream.Send(move)
	iter.OnLand(p)
}

func (p *Player) MoveTo(name string, pass bool) {
	iter := p.Pos.Next()
	total := 1
	for strings.Compare(iter.GetName(), name) != 0 {

		if pass {
			iter.OnPass(p)
		}

		iter = iter.Next()
		total++
	}

	p.Pos = iter
	iter.OnLand(p)
}

func (p *Player) BuyHousesFor(colour string, done bool) {

	if done {
		return
	}

	fmt.Printf("Buy a house for: \n")
	for i, prop := range p.Properties[colour] {
		if prop.Houses < 5 {
			fmt.Printf("%d - %s (%d)\n", i+1, prop.Position.Name, prop.HouseCost)
		}
	}
	fmt.Printf("E to escape\n")

	var choice string
	fmt.Scan(&choice)

	for !done {
		if strings.Compare(strings.ToLower(choice), "e") == 0 {
			return
		}

		index, err := strconv.Atoi(choice)

		if err != nil {
			fmt.Printf("%s is not a valid choice\n", choice)
			fmt.Scan(&choice)
			continue
		}

		if index-1 < 0 || index-1 >= len(p.Properties[colour]) {
			fmt.Printf("%s is not a valid choice\n", choice)
			fmt.Scan(&choice)
			continue
		}

		suc := p.Properties[colour][index].AddHouse()

		if !suc {
			fmt.Printf("%s Already has a hotel\n", p.Properties[colour][index].Position.Name)
		} else {
			fmt.Printf("Added a house to %s.\n", p.Properties[colour][index].Position.Name)
			fmt.Printf("%s has %d Houses\n", p.Properties[colour][index].Position.Name, p.Properties[colour][index].Houses)
		}

		defer p.BuyHousesFor(colour, false)
		return

	}
}

func (p *Player) BuyHouses(done bool) {

	if done {
		return
	}

	fmt.Printf("Buy houses for: \n")
	for i, colour := range p.FullSets {
		fmt.Printf("%d - %s\n", i+1, colour)
	}
	fmt.Printf("E to escape\n")

	var choice string
	fmt.Scan(&choice)
	for !done {

		if strings.Compare(strings.ToLower(choice), "e") == 0 {
			return
		}

		index, err := strconv.Atoi(choice)

		if err != nil {
			fmt.Printf("%s is not a valid choice\n", choice)
			fmt.Scan(&choice)
			continue
		}

		if index-1 < 0 || index-1 >= len(p.FullSets) {
			fmt.Printf("%s is not a valid choice\n", choice)
			fmt.Scan(&choice)
			continue
		}

		p.BuyHousesFor(p.FullSets[index-1], false)
		defer p.BuyHouses(false)
		return
	}
}

func (p *Player) SellProperty(c string) bool {
	fmt.Printf("Properties to sell ")

	count := 0
	if strings.Compare(c, "") == 0 {

		fmt.Printf(" (enter the colour to filter):\n")
		for colour, group := range p.Properties {
			fmt.Printf("%s\n", colour)
			for _, prop := range group {
				if prop.Houses == 0 {
					count++
					fmt.Printf("\t%d - %s (with Original price of %d)\n", count, prop.Position.Name, prop.Price)
				}
			}
		}
	} else {
		group, exists := p.Properties[c]
		if !exists {
			return false
		}

		fmt.Printf("in the colour %s:\n", c)
		for _, prop := range group {
			if prop.Houses == 0 {
				count++
				fmt.Printf("\t%d - %s (with Original price of %d)\n", count, prop.Position.Name, prop.Price)
			}
		}
	}

	fmt.Printf("E to escape\n")

	done := false
	var choice string

	for !done {

		fmt.Scan(&choice)

		//if they choose to exit
		if strings.Compare(strings.ToLower(choice), "e") == 0 {
			return true
		}

		_, cexs := p.Properties[choice]

		if cexs {
			p.SellProperty(choice)
			defer p.SellProperty("")
			return true
		}

		index, err := strconv.Atoi(choice)

		if err != nil {
			fmt.Printf("%s is an invalid choice\n", choice)
			continue
		}

		if (index < 1) || index > count {
			fmt.Printf("%s is an invalid choice\n", choice)
			continue
		}

		if strings.Compare("", c) == 0 {

			defer p.SellProperty("")
			innerCount := 0
			for g, group := range p.Properties {
				for prop_index, prop := range group {
					if prop.Houses == 0 {
						innerCount++
					}
					if innerCount == index {
						fmt.Printf("Are you sure you want to sell %s (Y/n)?\n", prop.Position.Name)
						var conf string
						fmt.Scan(&conf)
						if strings.Compare(strings.ToLower(conf), "y") == 0 {
							p.RemoveProperty(g, prop_index)
						}
						return true
					}
				}
			}
		} else {

			defer p.SellProperty(c)
			innerCount := 0
			for prop_index, prop := range p.Properties[c] {
				if prop.Houses == 0 {
					innerCount++
				}

				if innerCount == index {
					fmt.Printf("Are you sure you want to sell %s (Y/n)?\n", p.Properties[c][innerCount])
					var conf string
					fmt.Scan(&conf)
					if strings.Compare(strings.ToLower(conf), "y") == 0 {
						p.RemoveProperty(c, prop_index)
					}
					return true
				}
			}

		}
	}

	return false
}

func (p *Player) Sellhouses() {
	fmt.Printf("Sell houses from colour: \n")

	count := 0
	for _, group := range p.FullSets {
		for _, prop := range p.Properties[group] {
			if prop.Houses > 0 {
				count++
				fmt.Printf("\t%d - %s with %d houses (sell price %d)\n", count, prop.Position.Name, prop.Houses, prop.HouseCost/2)
			}
		}
	}

	fmt.Printf("E to escape/n")

	done := false

	var c string
	for !done {
		fmt.Scan(&c)

		if strings.Compare(strings.ToLower(c), "e") == 0 {
			return
		}

		index, err := strconv.Atoi(c)

		if err != nil {
			fmt.Printf("%s is a invalid choice\n", c)
			continue
		}

		if index < 0 || index > count {
			fmt.Printf("%s is a invalid choice\n", c)
			continue
		}

		innerCount := 0

		defer p.Sellhouses()

		for _, group := range p.FullSets {
			for _, prop := range p.Properties[group] {
				if prop.Houses > 0 {
					innerCount++
				}

				if innerCount == index {
					prop.RemoveHouse()
					return
				}
			}
		}
	}
}

func (p *Player) Mortgage() {
	fmt.Printf("Testing mortage\n")
}

func (p *Player) ListProperties() {

	if len(p.Properties) == 0 {
		fmt.Printf("No properties owned yes\n")
		return
	}

	for colour, group := range p.Properties {
		fmt.Printf("%s =================================\n", colour)
		if slices.Contains(p.FullSets, colour) {
			fmt.Printf("Full Set Owned\n")
		}
		for _, prop := range group {
			prop.Print()
		}
		fmt.Printf("====================================\n")
	}
}

func (p *Player) TurnInput(before bool) {
	fmt.Printf("What would you like to do this turn?\n")
	fmt.Printf("(L) - List properties.\n")
	fmt.Printf("(A) - List Available properties\n")
	fmt.Printf("(P) - Show Player profile\n")

	if before {
		fmt.Printf("(R) - Roll dice\n")
	} else {
		fmt.Printf("(M) - Manage assests\n")
		fmt.Printf("(E) - End turn\n")
	}

	done := false
	var c string
	for !done {
		fmt.Scan(&c)

		c = strings.ToLower(c)

		if (strings.Compare("r", c) == 0 && before) || (strings.Compare("e", c) == 0 && !before) {
			return
		}

		if strings.Compare("l", c) == 0 {
			p.ListProperties()
			defer p.TurnInput(before)
			return
		}

		if strings.Compare("a", c) == 0 {
			ListAvailable()
			defer p.TurnInput(before)
			return
		}

		if strings.Compare("m", c) == 0 && !before {
			p.ManageAssets()
			defer p.TurnInput(before)
			return
		}

		if strings.Compare("p", c) == 0 {
			p.Print()
			defer p.TurnInput(before)
			return
		}
		fmt.Printf("%s s an invalid choice\n", c)
	}
}

func (p *Player) ManageAssets() {
	if len(p.Properties) > 0 {
		fmt.Printf("How would you like to manage your assests?\n")
		if len(p.FullSets) > 0 {
			fmt.Printf("(B) Buy Houses\n")
			fmt.Printf("(SH) Sell Houses\n")
		}
		fmt.Printf("(S) Sell Property\n")
		fmt.Printf("(M) Mortgage Proprty\n")
		fmt.Printf("(E) Exit\n")

		var c string
		fmt.Scan(&c)

		c = strings.ToLower(c)

		if strings.Compare("e", c) == 0 {
			return
		}

		defer p.ManageAssets()

		// if buying houses
		if strings.Compare("b", c) == 0 {
			p.BuyHouses(false)
			return
		}

		if strings.Compare("sh", c) == 0 {
			p.Sellhouses()
			return
		}

		if strings.Compare("s", c) == 0 {
			p.SellProperty("")
			return
		}

		if strings.Compare("m", c) == 0 {
			p.Mortgage()
			return
		}
	} else {
		fmt.Printf("You do not have an assests to manage.\n")
		return
	}

}

func (p *Player) HandleJail() {
	t, double := RollDice()
	p.JailCount++

	if p.JailCount == 3 {
		p.SetFree()
		p.MovePlayer(true, 0)
		return
	}

	if p.JailFree > 0 {
		fmt.Printf("You have %d 'Get out of jail free cards' would you like to use one (Y/n):")
		var i string
		fmt.Scan(&i)
		if strings.Contains(strings.ToLower(i), "y") {
			p.JailFree = p.JailFree - 1
			p.SetFree()
			p.MovePlayer(true, 0)
			return
		}
	}

	if p.Cash >= 50 {
		fmt.Printf("You have %d to your name.\nWould you like to pay 50 to get out of jail (Y/n): ", p.Cash)

		var i string
		fmt.Scan(&i)

		if strings.Contains(strings.ToLower(i), "y") {
			p.Spend(50)
			p.SetFree()
			p.MovePlayer(true, 0)
			return
		}
	}

	if double {
		fmt.Printf("Rolled a double to get out of jail\n")
		p.SetFree()
		p.MovePlayer(false, t)
		return
	}

	fmt.Printf("Served %d rounds in jail\n", p.JailCount)
}

func (p *Player) MovePlayer(roll bool, count int) {
	goAgain := false
	if roll {
		// roll the dice
		t, double := RollDice()

		// if a double was rolled
		if double {
			// increase number of doubles
			p.DoubleCount = p.DoubleCount + 1
			goAgain = true
		} else {
			// resest double count
			p.DoubleCount = 0
		}

		count = t
	}

	// if three doubles have been thrown in a row
	if p.DoubleCount == 3 {
		// go to jail
		p.SetInJail()
		p.MoveTo("Jail", false)
		return
	}

	// Move through spaces
	for i := 0; i < count; i++ {
		if i != 0 {
			p.Pos.OnPass(p)
		}
		p.Pos = p.Pos.Next()
	}

	// print the Position landed on
	fmt.Printf("======= Position =======\n")
	p.Pos.Print()
	fmt.Printf("========================\n")

	p.Pos.OnLand(p)

	if goAgain && !p.Bankrupt {
		p.MovePlayer(true, 0)
	}
}
