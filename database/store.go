package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gianmarcomennecozzi/pro-haaukins/model"
	pb "github.com/gianmarcomennecozzi/pro-haaukins/proto"
	"log"
	"strconv"
	"sync"
	"time"
	"os"
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
	UpdateTeamSolvedChallenge(*pb.UpdateTeamSolvedChallengeRequest) (string, error)
	UpdateTeamLastAccess(*pb.UpdateTeamLastAccessRequest) (string, error)
	UpdateEventFinishDate(*pb.UpdateEventRequest) (string, error)
}

func NewStore() (Store, error){
	db, err := NewDBConnection()

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	err = InitTables(db)
	if err != nil {
		log.Printf("failed to init tables: %v", err)
	}
	return &store{ db: db }, nil
}

func NewDBConnection() (*sql.DB, error){

	host 		:= os.Getenv("DB_HOST")
	portString	:= os.Getenv("DB_PORT")
	dbUser     := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName     := os.Getenv("DB_NAME")

	//host 		:= "localhost"
	//portString	:= "5432"
	//dbUser     := "root"
	//dbPassword := "root"
	//dbName     := "mydb"

	port, err := strconv.Atoi(portString)
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, dbUser, dbPassword, dbName)
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

	addEventQuery := "INSERT INTO event (tag, name, available, capacity, frontends, exercises, started_at, finish_expected)" +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := s.db.Exec(addEventQuery, in.Tag, in.Name, in.Available, in.Capacity, in.Frontends, in.Exercises, in.StartTime, in.ExpectedFinishTime)

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
	addTeamQuery := "INSERT INTO team (id, event_tag, email, name, password, created_at, last_access, solved_challenges)" +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	_, err := s.db.Exec(addTeamQuery, in.Id, in.EventTag, in.Email, in.Name, in.Password, nowString, nowString, "[]")
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
		rows.Scan(&tag, &name, &available, &capacity, &frontends, &exercises, &startedAt, &expectedFinishTime, &finishedAt)
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

func (s store) UpdateTeamSolvedChallenge(in *pb.UpdateTeamSolvedChallengeRequest) (string, error) {
	s.m.Lock()
	defer s.m.Unlock()

	type Challenge struct {
		Tag  		string		`json:"tag"`
		CompletedAt string		`json:"completed-at"`
	}

	var solvedChallenges []Challenge
	var solvedChallengesDB string

	if err := s.db.QueryRow("SELECT solved_challenges FROM team WHERE id=$1", in.TeamId).Scan(&solvedChallengesDB); err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(solvedChallengesDB), &solvedChallenges); err != nil{
		return "", err
	}

	for _, sc := range solvedChallenges{
		if sc.Tag == in.Tag{
			return "", errors.New("challenge already solved")
		}
	}

	solvedChallenges = append(solvedChallenges, Challenge{
		Tag:         in.Tag,
		CompletedAt: in.CompletedAt,
	})

	newSolvedChallengesDB, _ := json.Marshal(solvedChallenges)

	_, err := s.db.Exec("UPDATE team SET solved_challenges = $2 WHERE id = $1", in.TeamId, string(newSolvedChallengesDB))
	if err != nil {
		return "", err
	}

	return "ok", nil
}

func (s store) UpdateTeamLastAccess(in *pb.UpdateTeamLastAccessRequest) (string, error) {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.db.Exec("UPDATE team SET last_access = $2 WHERE id = $1", in.TeamId, in.AccessAt)
	if err != nil {
		return "", err
	}

	return "ok", nil
}

func (s store) UpdateEventFinishDate(in *pb.UpdateEventRequest) (string, error) {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.db.Exec("UPDATE event SET finished_at = $2 WHERE id = $1", in.EventId, in.FinishedAt)
	if err != nil {
		return "", err
	}

	return "ok", nil
}