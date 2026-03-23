package database

import (
	// "database/sql"
	// "errors"
	// "fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	// "log"
	// "os"
	// "path/filepath"
	"time"
	"context"
)

func (db *appdbimpl) insertMessage(MessageBody string, UploaderId structs.Identifier, MessageId structs.Identifier, Date string, MediaType string, ConversationId structs.Identifier) error {
	_, err := db.c.ExecContext(context.Background(),
	`INSERT INTO message (MessageId, MessageBody,  Date, UploaderId, MediaType, ConversationId) VALUES (?, ?, ?, ?, ?, ?)`, 
	MessageId.Id, MessageBody, Date, UploaderId.Id, MediaType, ConversationId.Id)

	return err
}

func (db *appdbimpl) SendMessage(MessageBody string, UploaderId structs.Identifier, MediaType string, ConversationId structs.Identifier) (structs.Message, error) {

	thisMessageId := generateIdentifier("M")

	checkId, err := db.checkValidId(thisMessageId.Id, "M")

	if !checkId {
		return structs.Message{}, err
	}
	if err != nil {
		return structs.Message{}, err
	}

	messageDate := time.Now().UTC().Format(time.RFC3339)

	err = db.insertMessage(MessageBody, UploaderId, thisMessageId, messageDate, MediaType, ConversationId)
	if err != nil {
		return structs.Message{}, err
	}

	newMessage := structs.Message{
		MessageBody: MessageBody,
		Comments:    []structs.Comment{},
		UploaderId:  UploaderId,
		MessageId:   thisMessageId,
		Date:        messageDate,
		MediaType:   MediaType,
		ConversationId: ConversationId,
	}

	return newMessage, nil
}

func (db *appdbimpl) ForwardMessage(OldMessageId structs.Identifier, UploaderId structs.Identifier, ConversationId structs.Identifier) (structs.Message, error) {

	var MessageBody string
	var MediaType string

	thisMessageId := generateIdentifier("M")

	checkId, err := db.checkValidId(thisMessageId.Id, "M")

	if !checkId {
		return structs.Message{}, err
	}
	if err != nil {
		return structs.Message{}, err
	}

	messageDate := time.Now().UTC().Format(time.RFC3339)

	err = db.c.QueryRowContext(context.Background(),`SELECT MessageBody FROM message WHERE MessageId = ?`, OldMessageId).Scan(&MessageBody)
	if err != nil {
			return structs.Message{}, err
		}

	err = db.c.QueryRowContext(context.Background(),`SELECT MediaType FROM message WHERE MessageId = ?`, OldMessageId).Scan(&MediaType)
	if err != nil {
			return structs.Message{}, err
		}

	err = db.insertMessage(MessageBody, UploaderId, thisMessageId, messageDate, MediaType, ConversationId)
	if err != nil {
		return structs.Message{}, err
	}

	forwardedMessage := structs.Message{
		MessageBody: MessageBody,
		Comments:    []structs.Comment{},
		UploaderId:  UploaderId,
		MessageId:   thisMessageId,
		Date:        messageDate,
		MediaType:   MediaType,
		ConversationId: ConversationId,
	}

	return forwardedMessage, nil
}

func (db *appdbimpl) DeleteMessage(MessageId structs.Identifier) error {
	_, err := db.c.ExecContext(context.Background(),`DELETE FROM message WHERE MessageId = ?`, MessageId)
	return err
}
