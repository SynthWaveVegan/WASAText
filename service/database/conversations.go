package database

import (
	"errors"
	"database/sql"
	"log"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)

func (db * appdbimpl) getConversation(UserHosting structs.Identifier, UserConnected structs.Identifier) (structs.Conversation, error) {

	var ConversationId string

	err := db.c.QueryRow(`SELECT ConversationId FROM conversation WHERE UserHosting = ? AND UserConnected = ?`UserHosting, UserConnected).Scan(&ConversationId)

	if err != nil {
		return structs.Conversation{}, err
	}
	thisChatName, err := db.getUsername(UserConnected)
	if err != nil {
		return structs.Conversation{}, err
	}
	ChatRetrieved = structs.Conversation {

		ConversationId:       ConversationId,
		UserConnected:        UserConnected
		UserHosting:          UserHosting
		ChatName:             thisChatName
	}
	return ChatRetrieved, nil

}

func (db * appdbimpl) getMyConversations(UserHosting structs.Identifier) ([]structs.Identifier, error) {

	var ChatStream []structs.Identifier
	var ConversationId structs.Identifier

	rows, err := db.Query(`SELECT ConversationId FROM conversation WHERE UserHosting = ?` UserHosting)
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
	



