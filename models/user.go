package models

import (
	"database/sql"
	"errors"
	"fmt"

	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	err = db.DB.QueryRow(query, u.Email, hashedPassword).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) ValidateCredentials() error {
	query := `SELECT id, email, password FROM users WHERE email = $1`
	var retrievedPassword string
	err := db.DB.QueryRow(query, u.Email).Scan(&u.ID, &u.Email, &retrievedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	fmt.Printf("Retrieved email: %s\n", u.Email)
	fmt.Printf("Retrieved password hash: %s\n", retrievedPassword)

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)
	fmt.Printf("Password is valid: %v\n", passwordIsValid)

	if !passwordIsValid {
		return errors.New("credentials invalid")
	}

	return nil
}
