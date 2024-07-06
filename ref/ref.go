package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"termonopoly/termonopoly/comm"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MonopolyService struct {
	comm.UnimplementedTermonopolyServer
	players     []*Player
	startPlayer int
}

type Player struct {
	Name   string
	stream *comm.Termonopoly_GameStreamServer
}

func (s *MonopolyService) Broadcast(message comm.Message) {
	for _, pl := range s.players {
		(*pl.stream).Send(&message)
	}
}

func (s *MonopolyService) PickFirst() {
	for i := range s.players {
		j := rand.Intn(i + 1)
		s.players[i], s.players[j] = s.players[j], s.players[i]
	}
}

func (s *MonopolyService) GameStream(stream comm.Termonopoly_GameStreamServer) error {
	for {
		in, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		switch in.GetType() {

		case comm.MessageType_ACTION:
			switch in.GetAct().GetType() {
			case comm.ActionType_RECV:
				fmt.Printf("%s recieved some money")
				break
			case comm.ActionType_SEND:
				fmt.Printf("%s spent some money")
				break
			}
		case comm.MessageType_CONNECT:
			fmt.Printf("Connect request\n")

			fmt.Printf("New player %s\n", in.GetCon().GetPlayerName())

			if in.GetCon() == nil || strings.Compare("", in.GetCon().GetPlayerName()) == 0 {
				return status.Error(codes.InvalidArgument, "No Connect message supplied")
			}

			// add player to list of players
			s.players = append(s.players, &Player{
				stream: &stream,
				Name:   in.GetCon().GetPlayerName(),
			})

			message := &comm.Message{
				Type: comm.MessageType_CONNECT,
				Data: &comm.Message_Con{
					&comm.Connect{
						PlayerName: in.GetCon().GetPlayerName(),
					},
				},
			}
			stream.Send(message)

			if len(s.players) == 2 {

				message := &comm.Message{
					Type: comm.MessageType_START_TURN,
				}

				fmt.Printf("starting game\n")

				s.PickFirst()

				var order = make([]string, 0)
				for _, pl := range s.players {
					order = append(order, pl.Name)
				}

				for i, pl := range s.players {
					start_message := &comm.Message{
						Type: comm.MessageType_START_GAME,
						Data: &comm.Message_Con{
							&comm.Connect{
								Pos:   int32(i),
								Order: order,
							},
						},
					}
					(*pl.stream).Send(start_message)
				}

				(*s.players[s.startPlayer].stream).Send(message)
			}
			break

		case comm.MessageType_END_TURN:
			s.startPlayer = (s.startPlayer + 1) % 2

			(*s.players[s.startPlayer].stream).Send(&comm.Message{
				Type: comm.MessageType_START_TURN,
			})
			break
		default:
			fmt.Printf("Other command\n")
			break
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 5000))

	if err != nil {
		log.Fatal(err)
	}

	service := &MonopolyService{
		players:     make([]*Player, 0),
		startPlayer: 0,
	}

	grpcServer := grpc.NewServer()
	comm.RegisterTermonopolyServer(grpcServer, service)

	grpcServer.Serve(lis)
}
