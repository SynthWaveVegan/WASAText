package database

import (
	// "database/sql"
	// "errors"
	// "fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
)

// da capire con conversation id se si può direttamente usare il duo users
func (db *appdbimpl) getConversation(UserHosting structs.Identifier, UserConnected structs.Identifier) (structs.Conversation, error) {

	var thisUserHosting structs.Identifier
	var thisUserConnected structs.Identifier

	err := db.c.QueryRow(`SELECT (UserHosting, UserConnected) FROM conversation WHERE UserHosting = ? AND UserConnected = ?`, UserHosting, UserConnected).Scan(&thisUserHosting, &thisUserConnected)

	if err != nil {
		return structs.Conversation{}, err
	}
	thisChatName, err := db.getUsernamebyId(thisUserConnected)
	if err != nil {
		return structs.Conversation{}, err
	}
	ChatRetrieved := structs.Conversation{

		UserConnected: thisUserConnected,
		UserHosting:   thisUserHosting,
		ChatName:      thisChatName,
	}
	return ChatRetrieved, nil

}

func (db *appdbimpl) getMyConversations(UserHosting structs.Identifier) ([]structs.Identifier, error) {

	var ChatStream []structs.Identifier
	var UserConnected structs.Identifier

	rows, err := db.c.Query(`SELECT UserConnected FROM conversation WHERE UserHosting = ?`, UserHosting)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&UserConnected)
		if err != nil {
			return nil, err
		}
		ChatStream = append(ChatStream, UserConnected)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return ChatStream, nil

}
