package database

import (
	// "database/sql"
	// "errors"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	// "log"
	"math/rand"
	"time"
)

// startId: U(user) M(message) C(conversation) G(group) T(comment) P(photo) E(error)
func generateIdentifier(startId string) (structs.Identifier, error) {

	const lenght = 9

	rand.Seed(time.Now().UnixNano())

	var randomInt string
	for i := 0; i < 9; i++ {
		digit := rand.Intn(10)
		randomInt += fmt.Sprintf("%d", digit)
	}

	newId := structs.Identifier{
		Id: ("@" + startId + randomInt),
	}

	return newId, nil
}

func (db *appdbimpl) checkValidId(checkingId string, startId string) (bool, error) {

	var countCheck int
	var err error
	switch startId {
	case "U":
		err = db.c.QueryRow(`SELECT COUNT(*) FROM user WHERE UserId = ?`, checkingId).Scan(&countCheck)

	case "S":
		err = db.c.QueryRow(`SELECT COUNT(*) FROM conversation WHERE ConversationId = ?`, checkingId).Scan(&countCheck)

	case "M":
		err = db.c.QueryRow(`SELECT COUNT(*) FROM message WHERE MessageId = ?`, checkingId).Scan(&countCheck)

	case "G":
		err = db.c.QueryRow(`SELECT COUNT(*) FROM group WHERE GroupId = ?`, checkingId).Scan(&countCheck)

	case "C":
		err = db.c.QueryRow(`SELECT COUNT(*) FROM comment WHERE CommentId = ?`, checkingId).Scan(&countCheck)

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
	err := db.c.QueryRow(`SELECT Username FROM user WHERE UserId = ?`, UserId).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (db *appdbimpl) GetUserIdByName(Username string) (structs.Identifier, error) {
	var userId structs.Identifier
	err := db.c.QueryRow(`SELECT UserId FROM user WHERE Username = ?`, Username).Scan(&userId)
	if err != nil {
		return structs.Identifier{}, err
	}
	return userId, nil
}
