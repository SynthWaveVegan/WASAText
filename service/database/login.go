package database

import (
	// "database/sql"
	// "errors"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
)

func (db *appdbimpl) DoLogin(username string) (structs.Identifier, error) {

	var checkId bool
	var userId string

	
	log.Println("Verifica se l'utente esiste:", username)
	userId, checkId, err := db.CheckUserExist(username)
	if err != nil {
		log.Println("Errore nella verifica dell'utente:", err)
		return structs.Identifier{}, err
	}

	if checkId {
		log.Println("Id e Utente già esistenti", err)
		return structs.Identifier{Id: userId}, nil
	}

	
	newId := generateIdentifier("U")
	if err != nil {
		log.Println("Errore nel generare l'id:", err)
		return structs.Identifier{}, err
	}

	newUserId := newId.Id

	
	validId, err := db.checkValidId(newUserId, "U")
	if err != nil {
		log.Println("errore id non valido:", err)
		return structs.Identifier{}, err
	}

	
	if !validId {
		return structs.Identifier{}, fmt.Errorf("invalid user ID generated")
	}

	
	err = db.SetMyUsername("New", username, newUserId)
	if err != nil {
		log.Println("Errore nel set myusername:", err)
		return structs.Identifier{}, err
	}

	
	return structs.Identifier{Id: newUserId}, nil
}