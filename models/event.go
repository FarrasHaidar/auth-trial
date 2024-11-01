package models

import (
	"database/sql"
	"fmt"
	"time"

	"example.com/rest-api/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserId      int64
}

func (e *Event) Save() error {
	query := `
    INSERT INTO events(name, description, location, datetime, user_id) 
    VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := db.DB.QueryRow(query, e.Name, e.Description, e.Location, e.DateTime, e.UserId).Scan(&e.ID)
	return err
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT * FROM events`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserId)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := `SELECT id, name, description, location, datetime, user_id FROM events WHERE id = $1`

	var event Event
	err := db.DB.QueryRow(query, id).Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserId)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("event not found")
	} else if err != nil {
		return nil, err
	}

	return &event, nil
}

func (event *Event) Update() error {
	query := `
    UPDATE events SET name=$1, description=$2, location=$3, datetime=$4 WHERE id=$5`

	_, err := db.DB.Exec(query, event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}

func (event *Event) Delete() error {
	query := `DELETE FROM events WHERE id = $1`

	_, err := db.DB.Exec(query, event.ID)
	if err != nil {
		return fmt.Errorf("could not delete event: %w", err)
	}
	return err
}

func (e Event) Register(userId int64) error {
	query := `INSERT INTO registrations(event_id, user_id) VALUES ($1, $2)`

	_, err := db.DB.Exec(query, e.ID, userId)
	if err != nil {
		return fmt.Errorf("failed to register user with ID %d to event with ID %d: %w", userId, e.ID, err)
	}

	return nil
}

func (e Event) CancelRegistration(userId int64) error {
	query := `DELETE FROM registrations WHERE event_id = $1 AND user_id = $2`

	_, err := db.DB.Exec(query, e.ID, userId)
	if err != nil {
		return fmt.Errorf("failed to register user with ID %d to event with ID %d: %w", userId, e.ID, err)
	}

	return nil
}
