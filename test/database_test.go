package test

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
	"testing"
	"os"
)

func TestDBConnection(t *testing.T){
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
		t.Fatalf("Error: %s", err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlInfo)


	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
}
