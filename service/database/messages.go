package database

import (
	"errors"
	"database/sql"
	"log"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)

func (db * appdbimpl) sendMessage(MessageBody string Comments []Comment UploaderUserid Identifier MessageId Identifier Date string ConversationId Identifier) (error) {
	_, err := db.c.Exec("INSERT INTO message (MessageBody, Comments, UploaderUserid, MessageId, Date, ConversationId) VALUES (?, ?, ?, ?, ?, ?)" MessageBody, Comments, UploaderUserid, MessageId, Date, ConversationId)
	return err
}
func (db * appdbimpl) deleteMessage(messageId string) error {
	_, err := db.c.Exec("DELETE FROM message WHERE MessageId = ?", MessageId)
	return err
}