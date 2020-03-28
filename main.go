package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"context"
	pb "github.com/gianmarcomennecozzi/pro-haaukins/proto"
	"github.com/gianmarcomennecozzi/pro-haaukins/database"
)

const (
	port = ":50051"
)

type server struct {
	store database.Store
}


func (s server) AddEvent(context.Context, *pb.AddEventRequest) (*pb.InsertResponse, error) {
	panic("implement me")
}

func (s server) AddTeam(context.Context, *pb.AddTeamRequest) (*pb.InsertResponse, error) {
	panic("implement me")
}

func (s server) GetEvents(context.Context, *pb.Empty) (*pb.GetEventResponse, error) {
	result, err := s.store.GetEvents()
	if err != nil {
		return &pb.GetEventResponse{ErrorMessage: err.Error()}, nil
	}

	var events []*pb.GetEventResponse_Events
	for _, e := range result{
		events = append(events, &pb.GetEventResponse_Events{
			Name:                 e.Name,
			Tag:                  e.Tag,
			Frontends:            e.Frontends,
			Exercises:            e.Exercises,
			Available:            int32(e.Available),
			Capacity:             int32(e.Capacity),
			StartedAt:            e.StartedAt,
			ExpectedFinishTime:   e.ExpectedFinishTime,
			FinishedAt:           e.FinishedAt,
		})
	}

	return &pb.GetEventResponse{Events:events}, nil

}

func (s server) GetEventTeams(ctx context.Context,in *pb.GetEventTeamsRequest) (*pb.GetEventTeamsResponse, error) {
	result, err := s.store.GetTeams(in.EventTag)
	if err != nil {
		return &pb.GetEventTeamsResponse{ErrorMessage: err.Error()}, nil
	}

	var teams []*pb.GetEventTeamsResponse_Teams
	for _, t := range result{
		teams = append(teams, &pb.GetEventTeamsResponse_Teams{
			Id:                   t.Id,
			Email:                t.Email,
			Name:                 t.Name,
			CreatedAt:            t.CreatedAt,
			LastAccess:           t.LastAccess,
			SolvedChallenges:     t.SolvedChallenges,
		})
	}

	return &pb.GetEventTeamsResponse{Teams:teams}, nil
}

func main() {
	store, err := database.NewStore()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()


	pb.RegisterStoreServer(s, &server{store:store})
	fmt.Println("waiting client")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}



}