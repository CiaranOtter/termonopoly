package termonopoly

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
Structure of economics data:

0 - Nme of space
1 - Family (Colour group)
2 - Frequency (Not needed)
3 - Price (string with "$ in front"), if there is no cost of a chance it is blank
4 - Improvement cost
*/

func populate_space(space *Space, data *[][]string) (Node, error) {
	for index, field := range *data {

		// once the space type has been found
		if strings.Compare(field[0], space.Name) == 0 {

			if strings.Compare(field[1], "Corners") == 0 {
				// if it is a special corner

				*data = append((*data)[:index], (*data)[index+1:]...)

				if strings.Compare(field[0], "Go") == 0 {
					return &Go{
						Position: space,
					}, nil
				}

				if strings.Compare(field[0], "Jail") == 0 {
					return &Jail{
						Position: space,
					}, nil
				}

				if strings.Compare(field[0], "Go To Jail") == 0 {
					return &GoJail{
						Position: space,
					}, nil
				}

				if strings.Compare(field[0], "Free Parking") == 0 {
					return &Parking{
						Position: space,
					}, nil
				}
			}

			if strings.Compare(field[1], "Cards") == 0 {
				// if a card square
				*data = append((*data)[:index], (*data)[index+1:]...)

				if strings.Contains(field[0], "Chance") {
					return &Chance{
						Position: space,
					}, nil
				}

				if strings.Contains(field[0], "Community Chest") {
					return &Community{
						Position: space,
					}, nil
				}
			}

			if strings.Compare(field[1], "Tax") == 0 {
				// if a tax square

				*data = append((*data)[:index], (*data)[index+1:]...)
				return &Tax{
					Position: space,
					Price:    100,
				}, nil
			}

			// if utilities
			if strings.Compare(field[1], "Utilities") == 0 {

				v, err := strconv.Atoi(field[3])
				if err != nil {
					log.Fatal(err)
				}
				r, err := strconv.Atoi(field[5])
				if err != nil {
					log.Fatal(err)
				}

				return &Utility{
					Position: space,
					Price:    v,
					Rent:     r,
				}, nil
			}
			if strings.Compare(field[1], "Railroads") == 0 {

				r, err := strconv.Atoi(field[6])
				if err != nil {
					log.Fatal(err)
				}

				p, err := strconv.Atoi(field[3])

				if err != nil {
					log.Fatal(err)
				}

				return &Railroad{
					Position: space,
					Price:    p,
					Rent:     r,
				}, nil
				// if a railroad
			}
			// else it is a colour property

			prop := &Property{
				Position: space,
				Family:   field[1],
				Owner:    nil,
			}

			p, err := strconv.Atoi(field[3])

			if err != nil {
				log.Fatal(err)
			}

			prop.Price = p
			// set the price of a house
			hc, err := strconv.Atoi(field[4])
			if err != nil {
				log.Fatal(err)
			}

			prop.HouseCost = hc
			// set the base rent of the space
			br, err := strconv.Atoi(field[7])
			if err != nil {
				log.Fatal(err)
			}

			prop.Rent = br
			// set the rent for owning a whole set
			sr, err := strconv.Atoi(field[8])
			if err != nil {
				log.Fatal(err)
			}

			prop.SetRent = sr
			// init the array of rents based number of houses (hotel at index 5)
			prop.HouseRent = make([]int, 5)

			for i := 0; i < 5; i++ {

				r, err := strconv.Atoi(field[i+9])
				if err != nil {
					log.Fatal(err)
					prop.HouseRent[i] = r
				}
			}

			_, exists := Properties[prop.Family]

			if !exists {
				Properties[prop.Family] = make([]*Property, 0)
			}

			Properties[prop.Family] = append(Properties[prop.Family], prop)

			for _, t := range Properties[prop.Family] {
				t.TotalFamily = len(Properties[prop.Family])
			}

			return prop, nil

		}
	}

	return nil, errors.New("Failed to find the space")
}

func InitBoard() {
	file, err := os.Open("data/Space_Names.csv")

	if err != nil {
		fmt.Println(err)
	}

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	eco, err := os.Open("data/Economics.csv")

	if err != nil {
		fmt.Println(err)
	}

	reader = csv.NewReader(eco)
	eco_records, _ := reader.ReadAll()

	prev := Start

	Properties = make(map[string][]*Property)

	for index, line := range records {

		i, err := strconv.Atoi(line[0])

		if err != nil {
			log.Fatal(err)
		}

		space := &Space{
			Index: i,
			Name:  line[1],
			Next:  nil,
		}

		node, err := populate_space(space, &eco_records)

		if err != nil {
			log.Fatal(err)
		}

		if index == 0 {
			Start = node
			prev = Start
			continue
		}

		prev.SetNext(node)
		prev = prev.Next()

	}

	prev.SetNext(Start)

}

func GetAction(row []string) *ChanceAction {

	var action ChanceAction

	if strings.Compare(row[5], "adv") == 0 {
		action = &Advance{
			dest: row[6],
			pass: strings.Compare(row[6], "Jail") != 0,
		}
	}
	if strings.Compare(row[5], "sum") == 0 {
		action = &Sum{
			rollCount: 10,
		}
	}
	if strings.Compare(row[5], "near") == 0 {
		action = &Near{
			nodeType: row[6],
		}
	}
	if strings.Contains(row[5], "recv") {
		i, err := strconv.Atoi(row[6])

		if err != nil {
			log.Fatal("Failed to parse amount\n")
			return nil
		}

		if strings.Contains(row[5], "player") {
			action = &Trans{
				amount:    i,
				pay:       false,
				toPlayers: true,
			}
		} else {
			action = &Trans{
				amount:    i,
				pay:       false,
				toPlayers: false,
			}
		}

	}

	if strings.Compare(row[5], "jailFree") == 0 {
		action = &JailFree{}
	}

	if strings.Compare(row[5], "back") == 0 {

		dis, err := strconv.Atoi(row[6])

		if err != nil {
			log.Fatal("failed to parse distance")
			return nil
		}

		action = &Back{
			dist: dis,
		}
	}

	if strings.Compare(row[5], "repair") == 0 {

		h, err := strconv.Atoi(row[6])

		if err != nil {
			log.Fatal("failed to parse house price\n")
			return nil
		}

		hot, err := strconv.Atoi(row[7])

		if err != nil {
			log.Fatal("failed to parse hotel price\n")
			return nil
		}

		action = &Repair{
			house: h,
			hotel: hot,
		}
	}

	if strings.Contains(row[5], "pay") {
		am, err := strconv.Atoi(row[6])

		if err != nil {
			log.Fatal("Failed to parse amount")
			return nil
		}

		if strings.Contains(row[5], "player") {
			action = &Trans{
				amount:    am,
				pay:       true,
				toPlayers: true,
			}
		} else {
			action = &Trans{
				amount:    am,
				pay:       true,
				toPlayers: false,
			}
		}
	}

	return &action
}

func InitCards() {
	file, err := os.Open("data/Chance.csv")

	if err != nil {
		log.Fatal(err)
		return
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	ChanceCards = make([]*ChanceCard, 0)
	UsedChanceCards = make([]*ChanceCard, 0)

	for _, row := range records {
		card := ChanceCard{
			text:   row[1],
			action: *GetAction(row),
		}

		ChanceCards = append(ChanceCards, &card)
	}

	CommunityCards = make([]*CommunityCard, 0)
	UsedCommunityCards = make([]*CommunityCard, 0)

	comm_file, err := os.Open("data/CC.csv")

	if err != nil {
		log.Fatal(err)
	}

	reader = csv.NewReader(comm_file)
	com_records, _ := reader.ReadAll()

	for _, row := range com_records {
		card := CommunityCard{
			text:   row[1],
			action: *GetAction(row),
		}

		CommunityCards = append(CommunityCards, &card)
	}
}

func PrintBoard() {
	prev := Start

	for prev != nil {
		prev.Print()

		prev = prev.Next()
	}
}
