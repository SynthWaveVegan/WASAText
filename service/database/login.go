package database

import (
	"database/sql"
	"errors"
	// "fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
)

func (db *appdbimpl) DoLogin(username string) (structs.Identifier, error) {

	var checkId bool
	var userId string

	userId, checkId, err := db.CheckUserExist(username)

	if err != nil {
		return structs.Identifier{}, err
	}
	if checkId == true {
		return structs.Identifier{Id: userId}, nil
	}
	if checkId == false {

		newId, err := generateIdentifier("U")

		newUserId := newId.Id

		if err != nil {
			return structs.Identifier{}, err
		}

		validId, err := db.checkValidId(newUserId, "U")

		if validId == false {
			return structs.Identifier{}, err
		}
		if err != nil {
			return structs.Identifier{}, err
		}

		err = db.SetMyUsername("New", username, newUserId)
		if err != nil {
			return structs.Identifier{}, err
		}

		return structs.Identifier{Id: newUserId}, nil
	}
	return structs.Identifier{}, err
}