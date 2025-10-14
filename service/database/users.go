package database

import (
	"errors"
	"database/sql"
	"log"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)

func (db *appdbimpl) DoLogin(username string) (structs.Identifier, error) {

	//controllare se il nome esiste, se esiste ritorno quello altrimenti faccio setUsername, creo un nuovo id e lo ritorno
	var checkId bool
	var userId string

	userId, checkId, err = checkUserExist(username)

	if err != nil {
		return structs.Identifier{}, error
	}
	if checkId == true {
		return structs.Identifier{Id: userId}, nil
	}
	if checkId == false {

		err = db.SetMyUserame("New", username)

		if err != nil {
			return structs.Idenfifier{}, err
		}

		newId, err := generateIdentifier("U")

		newUserId = newId.Id

		if err != nil {
			return structs.Identifier{}, err
		}

		validId, err := db.checkValidId(newUserId, "U")

		if checkId == false {
			return structs.Identifier{}, err
		}
		if err != nil {
			return structs.Identifier{}, err
		}

		

		return structs.Identifier{Id: newUserId}, nil
	}
	
}
func (db *appdbimpl) SetMyUsername(mode string, newName string, userId string) error {

	var counter int
	validName, err = checkValidName(newName)

	err := db.c.QueryRow(`SELECT COUNT(*) FROM user WHERE Username = ?`, newName).Scan(&counter)
	if err != nil {
		return err
	}
	if counter != 0 {
		log.Printf("username taken")
		return nil
	}
	if validName == false {
		log.Printf("username invalid")
		return nil
	}
	if counter == 0 && validName == true {

		switch mode {

		case "New":

			err := db.createUser(newName, userId)
			return err

		case "Update":

			
			_, err := db.c.Exec(`UPDATE user SET Username = ? WHERE UserId = ?`, newName, userId)
			return err

		default:
			return err
		}
	}
	
	return nil
}

func (db * appdbimpl) CreateUser(username string, userId string) error {
	_, err := db.c.Exec(`INSERT INTO user (UserId, Username) VALUES (?, ?)` (userId, username))
	return err
}

func (db * appdbimpl) GetUser(UserId structs.Identifier) structs.User, error {
	var username string
	err := db.c.QueryRow(`SELECT Username FROM user WHERE UserId = ?`, UserId).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return structs.User{}, err
		} else {
			return structs.User{}, err
		}
	}
	
	UserFound = structs.User {
		Username: username
		UserId: UserId
	}

	return UserFound, nil


}

func (db * appdbimpl) checkUserExist(username string) (string, bool, error) {

	var userId string
	//questo controlla se esiste username nella table user e ritorna il corrispondente userId
	err := db.c.QueryRow(`SELECT UserId FROM user WHERE Username = ?`, username).Scan(&userid)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
				return false, "", nil
		} else {
			return false, "", err
		}
	} else {
		return true, userId, nil
	}

}

func (db * appdbimpl) getUsername(UserId structs.identifier) (string, error) {
	var username string
	err := db.c.QueryRow(`SELECT Username FROM user WHERE UserId = ?`, UserId).Scan(&username)
	if err != nil {
		return {}, err
	}
	return username, nil
}
		
