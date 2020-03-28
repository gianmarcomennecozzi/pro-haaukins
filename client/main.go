package main

import (
	"context"
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

	r, err := c.GetEventTeams(ctx, &pb.GetEventTeamsRequest{EventTag: "test"})
	if err != nil{
		log.Fatalf("could not greet: %v", err)
	}
	if r.ErrorMessage != ""{
		log.Fatalf("could not greet: %v", r.ErrorMessage)
	}
	for _, e := range r.Teams{
		log.Printf(e.Id)
	}
}