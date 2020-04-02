package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "github.com/gianmarcomennecozzi/pro-haaukins/proto"
)

const (
	address     = "localhost:50051"
)

func main(){

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewStoreClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//addEventRequest := pb.AddEventRequest{
	//	Name:                 "Test from Client",
	//	Tag:                  "clientestttttt",
	//	Frontends:            "awdwad,wadwad,rtr,trt",
	//	Exercises:            "bla,bla1,ciao",
	//	Available:            1212,
	//	Capacity:             20,
	//	ExpectedFinishTime:   "wadwad wdawadwadwa  awdadwad adwd",
	//}

	//addTeam := pb.AddTeamRequest{
	//	Id:                   "its_working",
	//	EventTag:             "menne",
	//	Email:                "menne@menne.com",
	//	Name:                 "menne",
	//	Password:             "menne_token_test",
	//}
	//r, err := c.AddTeam(ctx, &addTeam)

	//r, err := c.UpdateTeamSolvedChallenge(ctx, &pb.UpdateTeamSolvedChallengeRequest{
	//	TeamId:               "menne2",
	//	Tag:                  "prova",
	//	CompletedAt:          "prova time",
	//})
	r, err := c.GetEventTeams(ctx, &pb.GetEventTeamsRequest{EventTag:"merlo"})
	if err != nil{
		log.Fatalf("could not greet: %v", err)
	}
	if r.ErrorMessage != ""{
		log.Fatalf("my could not greet: %v", r.ErrorMessage)
	}
	//log.Println(r.Message)
	for _, e := range r.Teams{
		fmt.Println(e)
	}
}