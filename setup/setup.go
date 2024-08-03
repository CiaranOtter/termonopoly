package setup

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"termonopoly/game"
	"termonopoly/space"
)

func ReadCsv(file string, g *game.Game) game.SpaceInterface {
	f, err := os.Open(file)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()

	// fmt.Println(data)

	if err != nil {
		log.Fatal(err)
	}

	var prev game.SpaceInterface
	var first game.SpaceInterface

	for i := 1; i < len(data); i++ {
		row := data[i]

		t := space.SpaceFactory(row)

		if prev != nil {
			t.SetPrev(prev)
			prev.SetNext(t)
		} else {
			first = t
		}

		if strings.Compare(row[1], "Jail") == 0 {
			g.Jail = t
		}

		prev = t
	}

	prev.SetNext(first)
	first.SetPrev(prev)

	// fmt.Printf("\n")
	// for i := 0; i < 40; i++ {
	// 	first.Print()
	// 	first = first.GetNext()
	// }

	return first
}
