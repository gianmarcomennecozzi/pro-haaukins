package database

import "database/sql"

func InitTables(db *sql.DB) error{
	if _, err := createTables(db); err != nil {
		return err
	}
	if _, err := addEvent(db); err != nil {
		return err
	}
	if _, err := addTeam(db); err != nil {
		return err
	}
	return nil
}

func addEvent(db *sql.DB) (string, error){
	addEventQuery := "INSERT INTO event (tag, name, available, capacity, frontends, exercises, started_at, finish_expected)" +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	frontends := `[{"image": "kali","memoryMB": 1024,"cpu": 1}]`
	exercises := "bla,bladad,ciao"
	started_at := "2020-03-24T13:57:40.323156782+01:00"
	finish_expected := "2020-04-08T13:57:40.323156782+02:00"

	_, err := db.Exec(addEventQuery, "menne", "Menne Event Test", 5, 10, frontends, exercises, started_at, finish_expected)
	if err != nil {
		return "", err
	}
	return "ok", nil
}

func addTeam(db *sql.DB) (string, error){
	addTeamQuery := "INSERT INTO team (id, event_tag, email, name, password, created_at, last_access, solved_challenges)" +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	created_at := "2020-03-24T13:57:40.323156782+01:00"
	last_access := "2020-04-08T13:57:40.323156782+02:00"
	solvedChallenges := `[{"tag": "dwada","completed-at": "2020-03-24T14:23:21.102Z"},{"tag": "merlo","completed-at": "2020-03-24T14:23:21.102Z"}]`

	_, err := db.Exec(addTeamQuery, "its_working", "menne", "test@test.dk", "testttttt", "ciao", created_at, last_access, solvedChallenges)
	if err != nil {
		return "", err
	}
	return "ok", nil
}

func createTables(db *sql.DB) (string, error){

	//Create Event Table
	eventTableQuery := "create table Event(" +
		"tag varchar (50) primary key, " +
		"name varchar (150), " +
		"available integer, " +
		"capacity integer, " +
		"frontends text, " +
		"exercises text, " +
		"started_at varchar (100), " +
		"finish_expected varchar (100), " +
		"finished_at varchar (100));"
	if _ , err := db.Query(eventTableQuery); err != nil {
		return "", err
	}

	//Create Teams Table
	teamsTableQuery := "create table Team(" +
		"id varchar (50) primary key, " +
		"event_tag varchar (50), " +
		"email varchar (50), " +
		"name varchar (50), " +
		"password varchar (250), " +
		"created_at varchar (100), " +
		"last_access varchar (100), " +
		"solved_challenges text);"
	if _ , err := db.Query(teamsTableQuery); err != nil {
		return "", err
	}

	return "ok", nil
}