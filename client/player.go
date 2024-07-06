package main

import (
	"context"
	"fmt"
	"log"
	"termonopoly/termonopoly"
	"termonopoly/termonopoly/comm"

	"google.golang.org/grpc"
)

var Stream comm.Termonopoly_GameStreamClient
var MyPlayer termonopoly.Player

func startgame() {
	fmt.Printf("Starting the game\n")

	running := true

	for running {
		message, err := Stream.Recv()

		if err != nil {
			log.Fatal(err)
			Stream.CloseSend()
			running = false
			return
		}

		switch message.GetType() {
		case comm.MessageType_START_TURN:
			MyTurn()
			break
		}
	}
}

func MyTurn() {
	fmt.Printf("\nYour turn!\n")

	// ask player what they want to do
	MyPlayer.TurnInput(true)

	if MyPlayer.InJail {
		MyPlayer.HandleJail()
	} else {
		MyPlayer.MovePlayer(true, 0)
	}

	MyPlayer.TurnInput(false)

	message := &comm.Message{
		Type: comm.MessageType_END_TURN,
	}
	Stream.Send(message)
}

func main() {
	conn, err := grpc.NewClient("localhost:5000", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := comm.NewTermonopolyClient(conn)

	var name string
	fmt.Scan(&name)

	message := &comm.Message{
		Type: comm.MessageType_CONNECT,
		Data: &comm.Message_Con{
			&comm.Connect{
				PlayerName: name,
			},
		},
	}

	Stream, err = client.GameStream(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	Stream.Send(message)

	message, err = Stream.Recv()

	if err != nil {
		Stream.CloseSend()
		return
	}

	fmt.Printf("Connection successful\n")

	termonopoly.InitBoard()
	termonopoly.InitCards()
	termonopoly.Players = termonopoly.InitPlayers(2)

	waiting := true

	for waiting {
		fmt.Printf("waiting for the game to start...")
		start, err := Stream.Recv()

		if err != nil {
			log.Fatal(err)
			Stream.CloseSend()
		}

		switch start.GetType() {
		case comm.MessageType_PLAYER_JOIN:
			fmt.Printf("%s has joined the game\n")
			break

		case comm.MessageType_START_GAME:
			fmt.Printf("Starting the game\n")

			pos := start.GetCon().GetPos()
			order := start.GetCon().GetOrder()

			for i, name := range order {
				termonopoly.Players[i].Name = name
			}

			fmt.Printf("My position is %d\n", pos)

			MyPlayer = *termonopoly.Players[pos]
			MyPlayer.Stream = Stream

			waiting = false
			break

		default:
			break
		}
	}

	startgame()

	fmt.Printf("Game over\n")
}
