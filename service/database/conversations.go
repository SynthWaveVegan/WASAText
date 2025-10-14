package database

import (
	"errors"
	"database/sql"
	"log"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)

func (db * appdbimpl) getConversation() (structs.Conversation, error) {
	var  string

	err := db.c.QueryRow("SELECT  FROM conversation WHERE id=?").Scan(&)

	return 

}

func (db * appdbimpl) getMyConversations() ([]structs.Conversation, error) {
	

}

