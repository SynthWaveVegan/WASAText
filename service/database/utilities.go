package database

import (
	// "database/sql"
	// "errors"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
	"crypto/rand"
    "math/big"
	"context"
	
)

// startId: U(user) M(message) C(conversation) G(group) T(comment) P(photo) E(error)
func generateIdentifier(startId string) structs.Identifier {

	var randomInt string
	for i := 0; i < 9; i++ {
		digitBig, _ := rand.Int(rand.Reader, big.NewInt(10))
		digit := int(digitBig.Int64())
		randomInt += fmt.Sprintf("%d", digit)
	}

	newId := structs.Identifier{
		Id: ("@" + startId + randomInt),
	}

	return newId
}

func (db *appdbimpl) checkValidId(checkingId string, startId string) (bool, error) {
	var countCheck int
	var err error

	// Seleziona la query giusta in base al tipo di ID
	switch startId {
	case "U":
		err = db.c.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM users WHERE UserId = ?`, checkingId).Scan(&countCheck)
	case "C":
		err = db.c.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM conversation WHERE ConversationId = ?`, checkingId).Scan(&countCheck)
	case "M":
		err = db.c.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM message WHERE MessageId = ?`, checkingId).Scan(&countCheck)
	case "R":
		err = db.c.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM comment WHERE CommentId = ?`, checkingId).Scan(&countCheck)
	default:
		// Se startId non è valido, restituisci errore
		return false, fmt.Errorf("invalid startId value: %s", startId)
	}

	// Gestione degli errori generici
	if err != nil {
		return false, fmt.Errorf("database query failed: %w", err)
	}

	// Se il count è 0, l'ID è disponibile (non esiste nel DB)
	if countCheck == 0 {
		return true, nil
	}

	// Se l'ID esiste già, restituisci false
	return false, nil
}


func checkValidName(checkingName string) bool {
	
	log.Printf("Nome utente ricevuto: '%s'", checkingName)
	if len([]rune(checkingName)) >= 1 && len([]rune(checkingName)) <= 16 {
		return true
	}

	log.Printf("Invalid username length: '%s'. Length must be between 1 and 16 characters.", checkingName)
	return false
}

func (db *appdbimpl) getUsernamebyId(UserId structs.Identifier) (string, error) {
	var username string
	err := db.c.QueryRowContext(context.Background(),`SELECT Username FROM users WHERE UserId = ?`, UserId.Id).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (db *appdbimpl) GetUserIdByName(Username string) (structs.Identifier, error) {
	var userId string
	err := db.c.QueryRowContext(context.Background(),`SELECT UserId FROM users WHERE Username = ?`, Username).Scan(&userId)
	if err != nil {
		return structs.Identifier{}, err
	}

	UserId := structs.Identifier{
		Id: userId,
	}
	return UserId, nil
}
