package database

import (
	// "database/sql"
	// "errors"
	// "fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
	// "os"
	// "path/filepath"
	"time"
	"context"
)
func (db *appdbimpl) MarkMessageRead(UserId structs.Identifier, ConversationId structs.Identifier) error {

	var otherUserId string

	err := db.c.QueryRowContext(context.Background(), `
		SELECT u.UserId
		FROM users u
		JOIN userChat uc ON u.UserId = uc.UserId
		WHERE uc.ConversationId = ? AND u.UserId != ?`, ConversationId.Id, UserId.Id).Scan(&otherUserId)

	if err != nil {
		log.Println("ERROR QUERY (getting users):", err)
		return err
	}
	
	_, err = db.c.ExecContext(context.Background(), `UPDATE message SET IsRead = 'Yes' WHERE ConversationId = ? AND UploaderId = ? AND IsRead != 'Yes'`,
	 ConversationId.Id, otherUserId)
		if err != nil {
			return err
		}
	return nil
}

func (db *appdbimpl) insertMessage(MessageBody string, UploaderId structs.Identifier, MessageId structs.Identifier, Date string, MediaType string, ConversationId structs.Identifier, IsRead string, IsForwarded string) error {
	_, err := db.c.ExecContext(context.Background(),
	`INSERT INTO message (MessageId, MessageBody,  Date, UploaderId, MediaType, IsRead, IsForwarded, ConversationId) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 
	MessageId.Id, MessageBody, Date, UploaderId.Id, MediaType, IsRead, IsForwarded, ConversationId.Id)

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

	messageDate := time.Now().UTC().Format("02/01/2006 15:04:05")
	IsRead := "No" 
	IsForwarded := "No" 

	err = db.insertMessage(MessageBody, UploaderId, thisMessageId, messageDate, MediaType, ConversationId, IsRead, IsForwarded)
	if err != nil {
		return structs.Message{}, err
	}

	Uploader, err := db.GetUser(UploaderId)
	if err != nil {
		return structs.Message{}, err
	}

	newMessage := structs.Message{

		MessageBody: MessageBody,
		Comments:      []structs.Comment{},
		UploaderId:     UploaderId,
		Uploader:       Uploader,
		MessageId:      thisMessageId,
		Date:           messageDate,
		MediaType:      MediaType,
		ConversationId: ConversationId,
		IsRead:         IsRead,
		IsForwarded:    IsForwarded,
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

	messageDate := time.Now().UTC().Format("02/01/2006 15:04:05")
	IsRead := "No"
	IsForwarded := "Yes" 

	err = db.c.QueryRowContext(context.Background(),`SELECT MessageBody FROM message WHERE MessageId = ?`, OldMessageId.Id).Scan(&MessageBody)
	if err != nil {
			return structs.Message{}, err
		}

	err = db.c.QueryRowContext(context.Background(),`SELECT MediaType FROM message WHERE MessageId = ?`, OldMessageId.Id).Scan(&MediaType)
	if err != nil {
			return structs.Message{}, err
		}

	err = db.insertMessage(MessageBody, UploaderId, thisMessageId, messageDate, MediaType, ConversationId, IsRead, IsForwarded)
	if err != nil {
		return structs.Message{}, err
	}

	Uploader, err := db.GetUser(UploaderId)
	if err != nil {
		return structs.Message{}, err
	}

	forwardedMessage := structs.Message{
		MessageBody:    MessageBody,
		Comments:      []structs.Comment{},
		UploaderId:     UploaderId,
		Uploader:       Uploader,
		MessageId:      thisMessageId,
		Date:           messageDate,
		MediaType:      MediaType,
		ConversationId: ConversationId,
		IsRead:         IsRead,
		IsForwarded:    IsForwarded,
	}

	return forwardedMessage, nil
}

func (db *appdbimpl) DeleteMessage(MessageId structs.Identifier) error {
	_, err := db.c.ExecContext(context.Background(),`DELETE FROM message WHERE MessageId = ?`, MessageId.Id)
	return err
}
