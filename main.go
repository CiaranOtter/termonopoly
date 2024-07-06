package main

import (
	"fmt"
	"termonopoly/termonopoly"
)

var terminal bool

func IsGameOver(players *[]*termonopoly.Player) (bool, int) {
	count := 0
	win_index := -1

	for i, player := range *players {
		if player.IsBankrupt() {
			count++
		} else {
			win_index = i
		}
	}

	if count == len(*players)-1 {
		return true, win_index
	}

	return false, -1
}

func PrintWinner(winner *termonopoly.Player) {
	fmt.Printf("==============================\n")
	fmt.Printf("Game over\n")
	fmt.Printf("Winner is: \n\n")
	winner.Print()
}

func main() {
	termonopoly.InitBoard()
	termonopoly.InitCards()

	terminal = false
	player_count := 2
	termonopoly.Players = initPlayers(player_count)
	index := 0

	for {

		// check is game is over
		over, i := IsGameOver(&termonopoly.Players)

		// if the game is over print the winner and break the loop
		if over {
			PrintWinner(termonopoly.Players[i])
			break
		}

		// if the current index is bankrupt, skip their turn
		if termonopoly.Players[index].IsBankrupt() {
			index = (index + 1) % player_count
			continue
		}

		// print the player
		termonopoly.Players[index].Print()

		// Ask for player's before play move choice
		termonopoly.Players[index].TurnInput(true)

		if termonopoly.Players[index].InJail {
			// if the player is in jail
			termonopoly.Players[index].HandleJail()
		} else {
			// else play the players turn
			termonopoly.Players[index].MovePlayer(true, 0)
		}

		// ask for the players post move actions
		termonopoly.Players[index].TurnInput(false)

		// Move to the next player and continue
		index = (index + 1) % player_count
	}

	fmt.Printf("=============================\n")
}

func initPlayers(count int) []*termonopoly.Player {
	var players = make([]*termonopoly.Player, count)

	for i := 0; i < count; i++ {
		player := termonopoly.InitPlayer(fmt.Sprintf("Player %d", i))
		players[i] = &player
	}

	return players
}
