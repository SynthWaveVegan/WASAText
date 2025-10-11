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

		validId, err := checkValidId(newUserId, "U")

		if checkId == false {
			return structs.Identifier{}, err
		}
		if err != nil {
			return structs.Identifier{}, err
		}

		

		return structs.Identifier{Id: newUserId}, nil
	}
	
}
func (db *appdbimpl) SetMyUserame(mode string name string) error {

	







	_, err := db.c.Exec("INSERT INTO user (username) VALUES (?)", username)
	return err
}

checkUserExist(username string) (string, bool, error) {

	var userId string
	//questo controlla se esiste username nella table user e ritorna il corrispondente userId
	err := db.c.QueryRow(`SELECT userId FROM user WHERE username = ?`, username).Scan(&userid)

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
		
