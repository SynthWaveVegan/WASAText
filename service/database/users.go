package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
	"context"
)

func (db *appdbimpl) SetMyUsername(mode string, newName string, userId string) error {


	validName := checkValidName(newName)
	if !validName {
		log.Printf("username invalid")
		return fmt.Errorf("username '%s' is invalid", newName)  
	}

	
	var counter int
	err := db.c.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM users WHERE Username = ?`, newName).Scan(&counter)
	if err != nil {
		return fmt.Errorf("error checking username availability: %w", err) 
	}

	if counter != 0 {
		log.Printf("username '%s' is already taken", newName)
		return fmt.Errorf("username '%s' is already taken", newName) 
	}

	
	switch mode {

	case "New":
		
		err := db.CreateUser(newName, userId, "")
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err) 
		}

	case "Update":
		
		_, err := db.c.ExecContext(context.Background(), `UPDATE users SET Username = ? WHERE UserId = ?`, newName, userId)
		if err != nil {
			return fmt.Errorf("failed to update username: %w", err)  
		}

	default:
		return fmt.Errorf("invalid mode: %s", mode)  
	}

	return nil
}

func (db *appdbimpl) CreateUser(username string, userId string, userPhoto string) error {
	_, err := db.c.ExecContext(context.Background(),`INSERT INTO users (UserId, Username, UserPhoto) VALUES (?, ?, ?)`, userId, username, userPhoto)
	return err
}

func (db *appdbimpl) GetUser(UserId structs.Identifier) (structs.User, error) {

	var userphoto string
	var username string

	err := db.c.QueryRowContext(context.Background(),`SELECT Username FROM users WHERE UserId = ?`, UserId).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return structs.User{}, err
		} else {
			return structs.User{}, err
		}
	}

	err = db.c.QueryRowContext(context.Background(),`SELECT UserPhoto FROM users WHERE UserId = ?`, UserId).Scan(&userphoto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return structs.User{}, err
		} else {
			return structs.User{}, err
		}
	}
	

	UserFound := structs.User{
		Username: username,
		UserId:   UserId,
		UserPhoto: userphoto,
	}

	return UserFound, nil

}

func (db *appdbimpl) CheckUserExist(username string) (string, bool, error) {

	var userId string
	// Questo controlla se esiste l'username nella tabella 'users' e ritorna il corrispondente userId
	err := db.c.QueryRowContext(context.Background(), `SELECT UserId FROM users WHERE Username = ?`, username).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Se non c'è nessun risultato, significa che l'utente non esiste
			return "", false, nil
		}
		// Per altri errori (problemi di connessione, query malformate, etc.), restituisci un errore specifico
		return "", false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	// Se l'utente esiste, restituisci il suo userId
	return userId, true, nil
}

func (db *appdbimpl) SetMyPhoto(userId structs.Identifier, photoLink string) error {
	_, err := db.c.ExecContext(context.Background(),`UPDATE users SET UserPhoto = ? WHERE Userid = ?`, photoLink, userId)
	return err
}

