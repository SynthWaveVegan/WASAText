package database

import (
	"database/sql"
	"errors"
	// "fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
	"context"
)

func (db *appdbimpl) SetMyUsername(mode string, newName string, userId string) error {

	var counter int
	validName := checkValidName(newName)

	err := db.c.QueryRowContext(context.Background(),`SELECT COUNT(*) FROM user WHERE Username = ?`, newName).Scan(&counter)
	if err != nil {
		return err
	}
	if counter != 0 {
		log.Printf("username taken")
		return nil
	}
	if !validName {
		log.Printf("username invalid")
		return nil
	}
	if counter == 0 && validName {

		switch mode {

		case "New":

			err := db.CreateUser(newName, userId, "")
			return err

		case "Update":

			_, err := db.c.ExecContext(context.Background(),`UPDATE user SET Username = ? WHERE UserId = ?`, newName, userId)
			return err

		default:
			return err
		}
	}

	return nil
}

func (db *appdbimpl) CreateUser(username string, userId string, UserPhoto string) error {
	_, err := db.c.ExecContext(context.Background(),`INSERT INTO user (UserId, Username, UserPhoto) VALUES (?, ?, ?)`, userId, username, UserPhoto)
	return err
}

func (db *appdbimpl) GetUser(UserId structs.Identifier) (structs.User, error) {

	var userphoto string
	var username string

	err := db.c.QueryRowContext(context.Background(),`SELECT Username FROM user WHERE UserId = ?`, UserId).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return structs.User{}, err
		} else {
			return structs.User{}, err
		}
	}

	err = db.c.QueryRowContext(context.Background(),`SELECT UserPhoto FROM user WHERE UserId = ?`, UserId).Scan(&userphoto)
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
	//questo controlla se esiste username nella table user e ritorna il corrispondente userId
	err := db.c.QueryRowContext(context.Background(),`SELECT UserId FROM user WHERE Username = ?`, username).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, nil
		} else {
			return "", false, err
		}
	} else {
		return userId, true, nil
	}

}

func (db *appdbimpl) SetMyPhoto(userId structs.Identifier, photoLink string) error {
	_, err := db.c.ExecContext(context.Background(),`UPDATE user SET UserPhoto = ? WHERE Userid = ?`, photoLink, userId)
	return err
}

