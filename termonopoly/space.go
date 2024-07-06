package termonopoly

import (
	"fmt"
	"strings"
)

type Property struct {
	Position  *Space
	Price     int
	Family    string
	HouseCost int
	Rent      int
	SetRent   int
	HouseRent []int
	Owner     *Player

	Houses      int
	FullSet     bool
	TotalFamily int
}

func (p *Property) GetType() string {
	return "Property"
}

func (p *Property) AddHouse() bool {
	if p.Houses == 5 {
		return false
	}

	p.Houses++
	return true
}

func (p *Property) RemoveHouse() {
	p.Houses--
}

func (p *Property) GetName() string {
	return p.Position.Name
}

func (p *Property) Next() Node {
	return p.Position.Next
}

func (p *Property) SetNext(n Node) {
	p.Position.Next = n
}

func (s *Property) Print() {

	fmt.Printf("Name:\t\t%s\nType:\t\tProperty\n", s.Position.Name)
	fmt.Printf("Owned by:\t")
	if s.Owner == nil {
		fmt.Printf("No one")
	} else {
		fmt.Printf("%s", s.Owner.Name)
	}
	fmt.Printf("\n")

	fmt.Printf("Price:\t\t%d\nColour:\t\t%s\nTotal Colour:\t%d\nPrice of House:\t%d\nRent:\t\t%d\nFull Set Rent:\t%d\n", s.Price, s.Family, s.TotalFamily, s.HouseCost, s.Rent, s.SetRent)
}

func (s *Property) OnPass(player *Player) {

}

func (s *Property) OnLand(player *Player) {
	if s.Owner == nil {

		// if the player can afford the property
		if player.Cash >= s.Price {

			// offer it to them
			fmt.Printf("Buy %s for %d: ", s.Position.Name, s.Price)

			var i string
			fmt.Scan(&i)

			if strings.Contains(i, "Y") || strings.Contains(i, "y") {
				player.Spend(s.Price)
				player.AddProperty(s)
				s.Owner = player
				fmt.Printf("Property has been bought for %d by %s\n", s.Price, s.Owner.Name)
			}

		} else {
			fmt.Printf("You do not have enough cash to buy this property. Would you like to sell some assests?\n")
			var yn string
			fmt.Scan(&yn)

			if strings.Compare("y", strings.ToLower(yn)) == 0 {
				player.ManageAssets()
				defer s.OnLand(player)
				return
			}

		}

	} else if player != s.Owner {
		// if the space is owned by another player

		if player.Cash < s.Rent {
			var conf string

			bankrupt := false

			for !bankrupt {
				fmt.Printf("You can not afford %d rent with cash. Do you want to manage your assets?\n", s.Rent)
				fmt.Scan(&conf)
				if strings.Compare(strings.ToLower(conf), "y") == 0 {
					player.ManageAssets()
					defer s.OnLand(player)
					return
				} else {
					fmt.Printf("WARNING: you would be declaring bankruptcy. Are you sure you want to declare Bankruptcy?/n")
					var b string
					fmt.Scan(&b)

					if strings.Compare(strings.ToLower(b), "y") == 0 {
						fmt.Printf("Declaring bankruptcy\n")
						player.SetBankrupt()
						return
					}
				}
			}

		}

		player.Spend(s.Rent)
		s.Owner.Receive(s.Rent)
	}
}

type Railroad struct {
	Position *Space
	Price    int
	Rent     int
	Owner    *Player
}

func (p *Railroad) GetType() string {
	return "RailRoads"
}

func (p *Railroad) GetName() string {
	return p.Position.Name
}

func (p *Railroad) SetNext(n Node) {
	p.Position.Next = n
}

func (p *Railroad) Next() Node {
	return p.Position.Next
}

func (s *Railroad) Print() {
	fmt.Printf("Name:\t\t%s\nType:\t\tRail road\n", s.Position.Name)
	fmt.Printf("Owned by:\t")
	if s.Owner == nil {
		fmt.Printf("No one")
	} else {
		fmt.Printf("%s", s.Owner.Name)
	}
	fmt.Printf("\n")

	fmt.Printf("Price:\t\t%d\n", s.Price)
}

func (s *Railroad) OnPass(player *Player) {

}

func (s *Railroad) OnLand(player *Player) {

}

type Utility struct {
	Position *Space
	Price    int
	Rent     int
	Owner    *Player
}

func (p *Utility) GetType() string {
	return "Utility"
}

func (p *Utility) GetName() string {
	return p.Position.Name
}

func (p *Utility) SetNext(n Node) {
	p.Position.Next = n
}

func (p *Utility) Next() Node {
	return p.Position.Next
}

func (s *Utility) Print() {
	fmt.Printf("Name:\t\t%s\nType:\t\tUtility\n", s.Position.Name)
	fmt.Printf("Owned by:\t")
	if s.Owner == nil {
		fmt.Printf("No one")
	} else {
		fmt.Printf("%s", s.Owner.Name)
	}
	fmt.Printf("\n")

	fmt.Printf("Price:\t\t%d\n", s.Price)
}

func (s *Utility) OnPass(player *Player) {

}

func (s *Utility) OnLand(player *Player) {

}

type Jail struct {
	Position *Space
}

func (p *Jail) GetType() string {
	return "Jail"
}

func (p *Jail) GetName() string {
	return p.Position.Name
}

func (p *Jail) SetNext(n Node) {
	p.Position.Next = n
}

func (p *Jail) Next() Node {
	return p.Position.Next
}

func (j *Jail) Print() {
	fmt.Printf("Name:\t\t%s\nType:\t\tJail\n", j.Position.Name)
}

func (j *Jail) OnPass(player *Player) {

}

func (j *Jail) OnLand(player *Player) {

}

type Parking struct {
	Position *Space
}

func (p *Parking) GetType() string {
	return "Parking"
}

func (p *Parking) GetName() string {
	return p.Position.Name
}

func (p *Parking) SetNext(n Node) {
	p.Position.Next = n
}

func (p *Parking) Next() Node {
	return p.Position.Next
}

func (p *Parking) Print() {
	fmt.Printf("Name:\t\t%s\nType:\t\tParking\n", p.Position.Name)
}

func (p *Parking) OnPass(player *Player) {

}

func (p *Parking) OnLand(player *Player) {

}

type Go struct {
	Position *Space
}

func (p *Go) GetType() string {
	return "Begin"
}

func (p *Go) GetName() string {
	return p.Position.Name
}

func (p *Go) SetNext(n Node) {
	p.Position.Next = n
}

func (p *Go) Next() Node {
	return p.Position.Next
}

func (g *Go) Print() {
	fmt.Printf("Name:\t\t%s\nType:\t\tBegin\n", g.Position.Name)
}

func (g *Go) OnPass(player *Player) {
	player.Receive(200)
}

func (g *Go) OnLand(player *Player) {

}

type GoJail struct {
	Position *Space
}

func (p *GoJail) GetType() string {
	return "Go to jail"
}

func (p *GoJail) GetName() string {
	return p.Position.Name
}

func (p *GoJail) SetNext(n Node) {
	p.Position.Next = n
}

func (p *GoJail) Next() Node {
	return p.Position.Next
}
func (s *GoJail) Print() {
	fmt.Printf("Name:\t\t%s\nType:\t\tGo to jail\n", s.Position.Name)
}

func (s *GoJail) OnPass(player *Player) {

}

func (s *GoJail) OnLand(player *Player) {
	player.MoveTo("Jail", false)
}

type Tax struct {
	Position *Space
	Price    int
}

func (p *Tax) GetType() string {
	return "Tax"
}

func (p *Tax) GetName() string {
	return p.Position.Name
}

func (p *Tax) SetNext(n Node) {
	p.Position.Next = n
}

func (p *Tax) Next() Node {
	return p.Position.Next
}

func (t *Tax) Print() {
	fmt.Printf("Name:\t\t%s\nType:\t\tTax\n", t.Position.Name)
}

func (t *Tax) OnPass(player *Player) {

}

func (t *Tax) OnLand(player *Player) {

}

type Community struct {
	Position *Space
}

func (p *Community) GetType() string {
	return "Community Chest"
}

func (p *Community) GetName() string {
	return p.Position.Name
}

func (p *Community) SetNext(n Node) {
	p.Position.Next = n
}

func (p *Community) Next() Node {
	return p.Position.Next
}

func (c *Community) Print() {
	fmt.Printf("Name:\t\t%s\nType:\t\tCommunity chest\n", c.Position.Name)
}

func (c *Community) OnPass(player *Player) {

}

func (c *Community) OnLand(player *Player) {

}

type Chance struct {
	Position *Space
}

func (p *Chance) GetType() string {
	return "Chance"
}

func (p *Chance) GetName() string {
	return p.Position.Name
}

func (p *Chance) SetNext(n Node) {
	p.Position.Next = n
}

func (p *Chance) Next() Node {
	return p.Position.Next
}

func (c *Chance) Print() {
	fmt.Printf("Name:\t\t%s\nType:\t\tChance card\n", c.Position.Name)
}

func (c *Chance) OnPass(player *Player) {

}

func (c *Chance) OnLand(player *Player) {

	fmt.Printf("======= Chance card ========\n")
	card := DrawChanceCard()

	fmt.Printf("%s\n", card.text)

	card.action.OnDraw(player)
	fmt.Printf("============================\n")
}
