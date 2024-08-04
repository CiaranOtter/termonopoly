package space

import (
	"log"
	"strconv"
	"termonopoly/game"
	"termonopoly/styles"
)

var cols = map[string]styles.Colour{
	"Purple":     styles.ColourEnum.PURPLE,
	"Light Blue": styles.ColourEnum.CYAN,
	"Magenta":    styles.ColourEnum.MAGENTA,
	"Red":        styles.ColourEnum.RED,
	"Yellow":     styles.ColourEnum.YELLOW,
	"Green":      styles.ColourEnum.GREEN,
	"Dark Blue":  styles.ColourEnum.BLUE,
	"Orange":     styles.ColourEnum.WHITE,
}

type Property struct {
	BaseSpace

	// Name            string
	Colour styles.Colour
	// Price           int
	ImprovementCost int
	Rent            []int
	Owner           game.OwnerInterface
}

func (p *Property) SetOwner(owner game.OwnerInterface) {
	p.Owner = owner
}

func (p *Property) Afford(budget int) (int, bool) {
	if budget >= p.Price {
		return p.Price, true
	}

	return p.Price, false
}

func (p *Property) Print() {
	// fmt.Printf("Property: \033[%dm%s\033[0m\nColour: \033[%dm%d\033[0m\n", p.Colour, p.Name, p.Colour, p.Colour)
}

func (p *Property) OnLand(pl game.OwnerInterface) {
	// pl.HandleProperty(p)

	if p.Owner == nil {
		defer pl.OfferProperty(p)
		return
	}

	if pl != p.Owner {
		pl.ChargeRent(p.Rent[0])
	}
}

func NewProperty(row []string) *Property {
	// fmt.Printf("Property item\n")

	prop := &Property{
		Owner: nil,
	}
	prop.SetName(row[1])
	prop.SetPrice(row[3])
	// prop.SetRent(row[4])

	var err error
	rent := make([]int, 7)
	for i := 0; i < 7; i++ {
		rent[i], err = strconv.Atoi(row[7+i])
		if err != nil {
			log.Fatal(err)
		}
	}

	prop.Rent = rent
	prop.Colour = cols[row[2]]
	return prop
}
