package database

import (
	"errors"
	"database/sql"
	"log"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)


// startId: U(user) M(message) C(conversation) G(group) T(comment) P(photo) E(error)
func generateIdentifier(startId string) (structs,Identifier, error) {

	
	const lenght = 9

	rand.Seed(time.Now().UnixNano())

	var randomInt string
	for i := 0; i < 9; i++ {
		digit := rand.Intn(10) 
		randomInt += fmt.Sprintf("%d", digit)
	}

	newId := structs.Identifier {
		Identifier: ("@" + startId + randomInt)
	}
	
	return newId, nil
}

func checkValidId(checkingId string startId string) (bool, error) {

	var countCheck int 
	var err error
	switch startId {
		case "U":
			err = db.c.QueryRow(`SELECT COUNT(*) FROM user WHERE userId = ?`, checkingId).Scan(&countCheck)

		case "M":
			err = db.c.QueryRow(`SELECT COUNT(*) FROM message WHERE messageId = ?`, checkingId).Scan(&countCheck)

		case "C":
			err = db.c.QueryRow(`SELECT COUNT(*) FROM conversation WHERE conversationId = ?`, checkingId).Scan(&countCheck)

		case "G":
			err = db.c.QueryRow(`SELECT COUNT(*) FROM group WHERE groupId = ?`, checkingId).Scan(&countCheck)

		case "T":
			err = db.c.QueryRow(`SELECT COUNT(*) FROM comment WHERE commentId = ?`, checkingId).Scan(&countCheck)

		case "P":
			err = db.c.QueryRow(`SELECT COUNT(*) FROM photo WHERE photoId = ?`, checkingId).Scan(&countCheck)

		default: //da completare
	}

	if err != nil {
		return false, err
	}

	if countCheck == 0 {
		return true, nil
	}
	
	return false, err


}