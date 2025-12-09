package database

import (
	//"database/sql"
	//"errors"
	// "fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	//"log"
)

func (db *appdbimpl) DoLogin(username string) (structs.Identifier, error) {

	var checkId bool
	var userId string

	userId, checkId, err := db.CheckUserExist(username)

	if err != nil {
		return structs.Identifier{}, err
	}
	if checkId {
		return structs.Identifier{Id: userId}, nil
	}
	if !checkId {

		newId := generateIdentifier("U")

		newUserId := newId.Id

		if err != nil {
			return structs.Identifier{}, err
		}

		validId, err := db.checkValidId(newUserId, "U")

		if !validId {
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
