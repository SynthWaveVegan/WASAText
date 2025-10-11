package database

import (
	"errors"
	"database/sql"
	"log"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)

func (db * appdbimpl) getConversation() (structs.Conversation, error) {
	var emptyConversationId string
	err := db.c.QueryRow("SELECT conversationId FROM conversation WHERE id=?").Scan(&emptyConversationId)
	return 

}

func (db * appdbimpl) getMyConversations() ([]structs.Conversation, error) {
	

}

