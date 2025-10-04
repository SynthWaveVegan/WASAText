package database

import (
	"errors"
	"database/sql"
	"log"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)

func (db * appdbimpl) getConversation() (structs.Conversation, error) {

	var conversation structs.Conversation
	
	err := db.c.QueryRow("SELECT conversation FROM conversation WHERE id=?").Scan(&conversation.conversationId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	else return &conversation, nil
	
	

}

func (db * appdbimpl) getMyConversations() ([]structs.Conversation, error) {
	var stream []structs.Conversation
	//da completare
	err := db.c.QueryRow("SELECT * FROM conversation WHERE id=?")

}

