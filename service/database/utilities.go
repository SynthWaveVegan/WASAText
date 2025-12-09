package database

import (
	// "database/sql"
	// "errors"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	// "log"
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
	switch startId {
	case "U":
		err = db.c.QueryRowContext(context.Background(),`SELECT COUNT(*) FROM user WHERE UserId = ?`, checkingId).Scan(&countCheck)

	case "S":
		err = db.c.QueryRowContext(context.Background(),`SELECT COUNT(*) FROM conversation WHERE ConversationId = ?`, checkingId).Scan(&countCheck)

	case "M":
		err = db.c.QueryRowContext(context.Background(),`SELECT COUNT(*) FROM message WHERE MessageId = ?`, checkingId).Scan(&countCheck)

	case "G":
		err = db.c.QueryRowContext(context.Background(),`SELECT COUNT(*) FROM group WHERE GroupId = ?`, checkingId).Scan(&countCheck)

	case "C":
		err = db.c.QueryRowContext(context.Background(),`SELECT COUNT(*) FROM comment WHERE CommentId = ?`, checkingId).Scan(&countCheck)

	//case "P":
	//	err = db.c.QueryRow(`SELECT COUNT(*) FROM photo WHERE PhotoId = ?`, checkingId).Scan(&countCheck)

	default:
		return false, err
	}

	if err != nil {
		return false, err
	}

	if countCheck == 0 {
		return true, nil
	}

	return false, err

}

func checkValidName(checkingName string) bool {
	if len([]rune(checkingName)) <= 12 && len([]rune(checkingName)) >= 1 {
		return true
	}
	return false
}

func (db *appdbimpl) getUsernamebyId(UserId structs.Identifier) (string, error) {
	var username string
	err := db.c.QueryRowContext(context.Background(),`SELECT Username FROM user WHERE UserId = ?`, UserId).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (db *appdbimpl) GetUserIdByName(Username string) (structs.Identifier, error) {
	var userId structs.Identifier
	err := db.c.QueryRowContext(context.Background(),`SELECT UserId FROM user WHERE Username = ?`, Username).Scan(&userId)
	if err != nil {
		return structs.Identifier{}, err
	}
	return userId, nil
}
