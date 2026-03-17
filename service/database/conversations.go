package database

import (
	//"database/sql"
	//"errors"
	// "fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
	"context"
)

func (db *appdbimpl) CheckConversation(UserHosting string, UserConnected string) (structs.Identifier, bool, error) {

	var counter int
	err := db.c.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM conversation WHERE UserHosting = ? AND UserConnected = ?`, UserHosting, UserConnected).Scan(&counter)
	if err != nil {
		return structs.Identifier{}, false, err
	}
	if counter != 0 {
		var ConversationId structs.Identifier
		err := db.c.QueryRowContext(context.Background(),`SELECT ConversationId FROM conversation WHERE UserHosting = ? AND UserConnected = ?`, UserHosting, UserConnected).Scan(&ConversationId)
			if err != nil {
				return structs.Identifier{}, false, err
			}
		return ConversationId, false, nil
	}

	return structs.Identifier{}, true, nil
}
func (db *appdbimpl) CreateConversation(UserHosting structs.Identifier, UserConnectedName string) (structs.Conversation, error) {

	ConversationId := generateIdentifier("S")
	
	validId, err := db.checkValidId(ConversationId.Id, "S")
	if err != nil {
		return structs.Conversation{}, err
	}
	if !validId {
		return structs.Conversation{}, err
	}

	UserConnected, err := db.GetUserIdByName(UserConnectedName)
	if err != nil {
		return structs.Conversation{}, err
	}

	OtherConversationId, ChatNotExist, err := db.CheckConversation(UserHosting.Id, UserConnected.Id)
	if err != nil {
		return structs.Conversation{}, err
	}

	if ChatNotExist == false {
		OldChat, err := db.GetConversation(OtherConversationId)
		if err != nil {
			return structs.Conversation{}, err
		}

		return OldChat, nil
	}


	_, err = db.c.ExecContext(context.Background(),`INSERT INTO conversation (ConversationId, UserHosting, UserConnected) VALUES (?, ?, ?)`, ConversationId.Id, UserHosting.Id, UserConnected.Id)
	if err != nil {
		return structs.Conversation{}, err
	}

	NewChat := structs.Conversation{

		ConversationId: ConversationId,
		UserConnected:  UserConnected,
		UserHosting:    UserHosting,
		ChatName:       UserConnectedName,
	}

	return NewChat, nil
}

func (db *appdbimpl) GetConversation(ConversationId structs.Identifier) (structs.Conversation, error) {

	var thisUserHosting structs.Identifier
	var thisUserConnected structs.Identifier

	err := db.c.QueryRowContext(context.Background(),`SELECT (UserHosting, UserConnected) FROM conversation WHERE ConversationId = ?`, ConversationId).Scan(&thisUserHosting, &thisUserConnected)
	if err != nil {
		return structs.Conversation{}, err
	}

	thisChatName, err := db.getUsernamebyId(thisUserConnected)
	if err != nil {
		return structs.Conversation{}, err
	}
	ChatRetrieved := structs.Conversation{

		ConversationId: ConversationId,
		UserConnected:  thisUserConnected,
		UserHosting:    thisUserHosting,
		ChatName:       thisChatName,
	}

	return ChatRetrieved, nil

}

func (db *appdbimpl) GetMyConversations(UserHosting structs.Identifier) ([]structs.Identifier, error) {

	var ChatStream []structs.Identifier
	var ConversationId structs.Identifier

	rows, err := db.c.QueryContext(context.Background(),`SELECT ConversationId FROM conversation WHERE UserHosting = ?`, UserHosting.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&ConversationId.Id)
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
