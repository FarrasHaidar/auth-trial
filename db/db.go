package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	connStr := "user=postgres password=Xerox247@ dbname=postgres host=localhost port=5432 sslmode=disable"
	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Could not ping the database:", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	fmt.Println("Connected to PostgreSQL successfully")

	createTables()
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
	)`

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Could not create table:", err)
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime TIMESTAMP NOT NULL,
		user_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(createEventsTable)
	if err != nil {
		log.Fatal("Could not create table:", err)
	}
	fmt.Println("Table created successfully")

	createRegistrationTable := `
	CREATE TABLE IF NOT EXISTS registrations (
	    id SERIAL PRIMARY KEY,
        event_id INTEGER,
        user_id INTEGER,
        FOREIGN KEY (event_id) REFERENCES events(id),
        FOREIGN KEY (user_id) REFERENCES users(id))
	`
	_, err = DB.Exec(createRegistrationTable)

	if err != nil {
		log.Fatal("Could not create table:", err)
	}
}
