package database

import (
	"database/sql"
	"errors"
	// "fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
)

func (db * appdbimpl) createConversation(UserHosting structs.Identifier, UserConnected structs.Identifier) (structs.Conversation, error) {

	ConversationId, err := generateIdentifier("S")
	if err != nil {
		return structs.Conversation{}, err
	}
	validId, err := db.checkValidId(ConversationId.Id, "S")
	if err != nil {
		return structs.Conversation{}, err
	}
	if validId == false {
		return structs.Conversation{}, err
	}

	var ChatName string

	err = db.c.Query(`SELECT Username FROM user WHERE UserId = ?`, UserConnected).Scan(&ChatName)
	if err != nil {
		return structs.Conversation{}, err
	}

	_, err := db.c.Exec(`INSERT INTO conversation (ConversationId, UserHosting, UserConnected) VALUES (?, ?, ?)`, ConversationId.Id, UserHosting, UserConnected)
	if err != nil {
		return structs.Conversation{}, err
	}

	NewChat := structs.Conversation{

		ConversationId: ConversationId.Id,
		UserConnected:  UserConnected,
		UserHosting:    UserHosting,
		ChatName:       ChatName,
	}

	return NewChat, nil
}

func (db *appdbimpl) getConversation(ConversationId structs.Identifier) (structs.Conversation, error) {

	var thisUserHosting structs.Identifier
	var thisUserConnected structs.Identifier

	err := db.c.QueryRow(`SELECT (UserHosting, UserConnected) FROM conversation WHERE ConversationId = ?`, ConversationId).Scan(&thisUserHosting, &thisUserConnected)
	if err != nil {
		return structs.Conversation{}, err
	}

	thisChatName, err := db.getUsernamebyId(thisUserConnected)
	if err != nil {
		return structs.Conversation{}, err
	}
	ChatRetrieved := structs.Conversation{

		ConversationId: ConversationId,
		UserConnected: thisUserConnected,
		UserHosting:   thisUserHosting,
		ChatName:      thisChatName,
	}

	return ChatRetrieved, nil

}

func (db *appdbimpl) getMyConversations(UserHosting structs.Identifier) ([]structs.Identifier, error) {

	var ChatStream []structs.Identifier
	var ConversationId structs.Identifier

	rows, err := db.c.Query(`SELECT ConversationId FROM conversation WHERE UserHosting = ?`, UserHosting)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&ConversationId)
		if err != nil {
			return nil, err
		}
		ChatStream = append(ChatStream, ConversationId)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return ChatStream, nil

}
