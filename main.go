package main

import (
	"fmt"
	"math/rand"
	"termonopoly/game"
	"termonopoly/player"
	"termonopoly/setup"
	"time"
)

func main() {
	game := game.NewGame()

	start := setup.ReadCsv("./data/Economics.csv", game)

	pl := player.NewPlayer("Ciaran", start, 1500)

	pl.Space.Print()

	for {
		roll := rand.Intn(12)
		fmt.Printf("Moved forward %d spaces\n", roll)
		pl.Move(roll, player.FORWARD)
		time.Sleep(1 * time.Second)
	}

}
