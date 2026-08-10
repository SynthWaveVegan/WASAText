package database

import (
	// "database/sql"
	// "errors"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
	// "os"
	// "path/filepath"
	"time"
	"context"
)
func (db *appdbimpl) MarkMessageRead(UserId structs.Identifier, ConversationId structs.Identifier) error {

	var IsGroup string
	err := db.c.QueryRowContext(context.Background(), `SELECT IsGroup FROM conversation WHERE ConversationId = ?`, ConversationId.Id).Scan(&IsGroup)
	if err != nil {
		return err
	}
	
	if IsGroup == "no" {

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
	
		_, err = db.c.ExecContext(context.Background(), `UPDATE message SET IsRead = 'yes' WHERE ConversationId = ? AND UploaderId = ? AND IsRead != 'yes'`, ConversationId.Id, otherUserId)
			if err != nil {
				return err
			}
	} else {

		_, err = db.c.ExecContext(context.Background(), `UPDATE message SET IsRead = 'yes' WHERE ConversationId = ? AND UploaderId = ? AND IsRead != 'yes'`, ConversationId.Id, UserId.Id)
		if err != nil {
			return err
		}

		var unreadCount int

		err = db.c.QueryRowContext(
			context.Background(),
			`
			SELECT COUNT(*)
			FROM message m
			JOIN userChat uc ON m.ConversationId = uc.ConversationId
			WHERE m.ConversationId = ?
				AND m.IsRead != 'yes'
				AND uc.UserId != ?
			`,
			ConversationId.Id,
			UserId.Id,
		).Scan(&unreadCount)

		if err != nil {
			return err
		}
		

		if unreadCount == 0 {
			_, err = db.c.ExecContext(context.Background(), `
				UPDATE message
				SET IsRead = 'yes'
				WHERE ConversationId = ?
				AND IsRead != 'yes'`,
				ConversationId.Id)
			if err != nil {
				return err
			}
		}

		
	}

	return nil
}

func (db *appdbimpl) insertMessage(MessageBody string, UploaderId structs.Identifier, MessageId structs.Identifier, Date string, MediaType string, ConversationId structs.Identifier, IsRead string, IsForwarded string) error {
	_, err := db.c.ExecContext(context.Background(), `INSERT INTO message (MessageId, MessageBody,  Date, UploaderId, MediaType, IsRead, IsForwarded, ConversationId) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, MessageId.Id, MessageBody, Date, UploaderId.Id, MediaType, IsRead, IsForwarded, ConversationId.Id)

	return err
}

func (db *appdbimpl) SendMessage(MessageBody string, UploaderId structs.Identifier, MediaType string, ConversationId structs.Identifier) (structs.Message, error) {

	thisMessageId := generateIdentifier("M")

	checkId, err := db.checkValidId(thisMessageId.Id, "M")

	if !checkId {
		return structs.Message{}, fmt.Errorf("invalid id")
	}
	if err != nil {
		return structs.Message{}, err
	}

	messageDate := time.Now().UTC().Format("02/01/2006 15:04:05")
	IsRead := "no" 
	IsForwarded := "no" 

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
		return structs.Message{}, fmt.Errorf("invalid id")
	}
	if err != nil {
		return structs.Message{}, err
	}

	messageDate := time.Now().UTC().Format("02/01/2006 15:04:05")
	IsRead := "no"
	IsForwarded := "yes" 

	err = db.c.QueryRowContext(context.Background(), `SELECT MessageBody, MediaType FROM message WHERE MessageId = ?`, OldMessageId.Id,).Scan(&MessageBody, &MediaType)
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
