package database

import (
	"database/sql"
	"fmt"
	"log"
	pb "pro-haaukins/proto"
	"sync"
	"time"
	"pro-haaukins/model"
)
const (
	HOST		= "127.0.0.1"
	POST		= 5432
	DB_USER     = "root"
	DB_PASSWORD = "root"
	DB_NAME     = "mydb"
)

type store struct {
	m sync.Mutex
	db *sql.DB
}

type Store interface {
	AddEvent(*pb.AddEventRequest) (string, error)
	AddTeam(*pb.AddTeamRequest) (string, error)
	GetEvents() ([]model.Event, error)
	GetTeams(string) ([]model.Team, error)
}

func NewStore() (Store, error){
	db, err := NewDBConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return &store{ db: db }, nil
}

func NewDBConnection() (*sql.DB, error){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		HOST, POST, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", psqlInfo)


	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (s store) AddEvent(in *pb.AddEventRequest) (string, error) {
	s.m.Lock()
	defer s.m.Unlock()

	now := time.Now()
	nowString := now.Format("2006-01-02 15:04:05")
	addEventQuery := "INSERT INTO event (tag, name, available, capacity, frontends, exercises, started_at, finish_expected)" +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := s.db.Exec(addEventQuery, in.Tag, in.Name, in.Available, in.Capacity, in.Frontends, in.Exercises, nowString, in.ExpectedFinishTime)

	if err != nil {
		return "", err
	}
	return "Event correctly added!", nil
}

func (s store) AddTeam(in *pb.AddTeamRequest) (string, error){
	s.m.Lock()
	defer s.m.Unlock()

	now := time.Now()
	nowString := now.Format("2006-01-02 15:04:05")
	addTeamQuery := "INSERT INTO team (id, event_tag, email, name, password, created_at, last_access)" +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	_, err := s.db.Exec(addTeamQuery, in.Id, in.EventTag, in.Email, in.Name, in.Password, nowString, nowString)
	if err != nil {
		return "", err
	}
	return "Team correctly added!", nil
}

func (s store) GetEvents() ([]model.Event, error) {

	s.m.Lock()
	defer s.m.Unlock()

	rows, err := s.db.Query("SELECT * FROM event")
	if err != nil{
		return []model.Event{}, err
	}
	var events []model.Event
	for rows.Next() {
		var tag 				string
		var name 				string
		var frontends 			string
		var exercises 			string
		var available 			uint
		var capacity 			uint
		var startedAt			string
		var expectedFinishTime 	string
		var finishedAt			string
		rows.Scan(&tag, &name, &frontends, &exercises, &available, &capacity, &startedAt, &expectedFinishTime, &finishedAt)
		events = append(events, model.Event{
			Tag:                tag,
			Name:               name,
			Frontends:          frontends,
			Exercises:          exercises,
			Available:          available,
			Capacity:           capacity,
			StartedAt:          startedAt,
			ExpectedFinishTime: expectedFinishTime,
			FinishedAt:         finishedAt,
		})
	}

	return events, nil
}

func (s store) GetTeams(tag string) ([]model.Team, error) {
	s.m.Lock()
	defer s.m.Unlock()

	rows, err := s.db.Query("SELECT * FROM team WHERE event_tag=$1", tag)
	if err != nil{
		return []model.Team{}, err
	}

	var teams []model.Team
	for rows.Next(){
		var id 					string
		var eventTag			string
		var email				string
		var name				string
		var password			string
		var createdAt			string
		var lastAccess			string
		var solvedChallenges	string

		rows.Scan(&id, &eventTag, &email, &name, &password, &createdAt, &lastAccess, &solvedChallenges)
		teams = append(teams, model.Team{
			Id:               id,
			EventTag:         eventTag,
			Email:            email,
			Name:             name,
			Password:         password,
			CreatedAt:        createdAt,
			LastAccess:       lastAccess,
			SolvedChallenges: solvedChallenges,
		})
	}
	return teams, nil
}
