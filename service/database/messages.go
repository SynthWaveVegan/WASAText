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
)

func (db *appdbimpl) insertMessage(MessageBody string, UploaderId structs.Identifier, MessageId structs.Identifier, Date string) error {
	_, err := db.c.Exec(`INSERT INTO message (MessageId, MessageBody,  Date, UploaderId) VALUES (?, ?, ?, ?)`, MessageId, MessageBody, Date, UploaderId)

	return err
}

func (db *appdbimpl) sendMessage(MessageBody string, UploaderId structs.Identifier) (structs.Message, error) {

	var thisMessageId structs.Identifier
	var checkId = false

	thisMessageId, err := generateIdentifier("M")
	if err != nil {
		return structs.Message{}, err
	}

	checkId, err = db.checkValidId(thisMessageId.Id, "M")

	if checkId == false {
		return structs.Message{}, err
	}
	if err != nil {
		return structs.Message{}, err
	}

	messageDate := time.Now().UTC().Format(time.RFC3339)

	err = db.insertMessage(MessageBody, UploaderId, thisMessageId, messageDate)
	if err != nil {
		return structs.Message{}, err
	}

	newMessage := structs.Message{
		MessageBody: MessageBody,
		Comments:    []structs.Comment{},
		UploaderId:  UploaderId,
		MessageId:   thisMessageId,
		Date:        messageDate,
	}

	return newMessage, nil
}

func (db *appdbimpl) forwardMessage(OldMessageId structs.Identifier, UploaderId structs.Identifier) (structs.Message, error) {

	var MessageBody string

	thisMessageId, err := generateIdentifier("M")
	if err != nil {
		return structs.Message{}, err
	}

	checkId, err := db.checkValidId(thisMessageId.Id, "M")

	if checkId == false {
		return structs.Message{}, err
	}
	if err != nil {
		return structs.Message{}, err
	}

	messageDate := time.Now().UTC().Format(time.RFC3339)

	err = db.c.QueryRow(`SELECT MessageBody FROM message WHERE MessageId = ?`, OldMessageId).Scan(&MessageBody)

	err = db.insertMessage(MessageBody, UploaderId, thisMessageId, messageDate)
	if err != nil {
		return structs.Message{}, err
	}

	forwardedMessage := structs.Message{
		MessageBody: MessageBody,
		Comments:    []structs.Comment{},
		UploaderId:  UploaderId,
		MessageId:   thisMessageId,
		Date:        messageDate,
	}

	return forwardedMessage, nil
}

func (db *appdbimpl) deleteMessage(MessageId structs.Identifier) error {
	_, err := db.c.Exec(`DELETE FROM message WHERE MessageId = ?`, MessageId)
	return err
}
