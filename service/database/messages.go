package database

import (
	"errors"
	"database/sql"
	"log"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)

func (db * appdbimpl) insertMessage(MessageBody string, Comments []structs.Comment, UploaderUserid structs.Identifier, MessageId structs.Identifier, Date string, ConversationId structs.Identifier) (error) {
	_, err := db.c.Exec(`INSERT INTO message (MessageBody, Comments, Uploaderid, MessageId, Date, ConversationId) VALUES (?, ?, ?, ?, ?, ?)` (MessageBody, Comments, Uploaderid, MessageId, Date, ConversationId))

	return err
}

func (db * appdbimpl) sendMessage(MessageBody string, UploaderId structs.Identifier, ConversationId structs.Identifier) (structs.Message , error) {

	var thisMessageId structs.Identifier
	var checkId false

	thisMessageId, err = generateIdentifier("M")
	if err != nil {
		return nil, err
	}
	
	checkId, err = checkValidId(thisMessageId, "M")

	if checkId == false {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	messageDate := time.Now().UTC().Format(time.RFC3339)

	err := db.insertMessage(MessageBody, []structs.Comment{}, UploaderId, thisMessageId, messageDate, ConversationId)
	if err != nil {
		return nil, err
	}


	newMessage = structs.Message {
		MessageBody:          MessageBody,
		Comments:             []structs.Comment{},
		Uploaderid:           UploaderId,
		MessageId:            thisMessageId,
		Date:                 messageDate,
		ConversationId:       ConversationId,
	}

	return newMessage, nil
}

func (db * appdbimpl) forwardMessage(ConversationId structs.Identifier, OldMessageId structs.Identifier, UploaderId structs.Identifier) (structs.Message, error) {

	var MessageBody string

	thisMessageId, err = generateIdentifier("M")
	if err != nil {
		return nil, err
	}
	
	checkId, err = checkValidId(thisMessageId, "M")

	if checkId == false {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	messageDate := time.Now().UTC().Format(time.RFC3339)
	
	err := db.c.QueryRow(`SELECT MessageBody FROM message WHERE MessageId = ?`, OldMessageId).Scan(&MessageBody)


	err := db.insertMessage(MessageBody, []structs.Comment{}, UploaderId, thisMessageId, messageDate, ConversationId)
	if err != nil {
		return nil, err
	}


	forwardedMessage = structs.Message {
		MessageBody:          MessageBody,
		Comments:             []structs.Comment{},
		Uploaderid:           UploaderId,
		MessageId:            thisMessageId,
		Date:                 messageDate,
		ConversationId:       ConversationId,
	}

	return forwardedMessage, nil

}






func (db * appdbimpl) deleteMessage(MessageId structs.Identifier) error {
	_, err := db.c.Exec(`DELETE FROM message WHERE MessageId = ?`, MessageId)
	return err
}

func (db * appdbimpl) insertComment(CommentId structs.Identifier, MessageId structs.Identifier, CommentBody string, Date string) error {

	_, err := db.c.Exec(`INSERT INTO comment(CommentId, MessageId, CommentBody, Date) VALUES (?, ?, ?, ?)` (CommentId, MessageId, CommentBody, Date))
	
	return err
}

func (db * appdbimpl) commentMessage(CommentBody string, MessageId structs.Identifier) (structs.Comment, error) {

	var thisCommentId string
	var checkId false

	thisCommentId, err = generateIdentifier("T")
	if err != nil {
		return nil, err
	}
	checkId, err = checkValidId(thisCommentId, "T")
	if checkId == false {
		return nil, err
	}
	if err != nil {
		return nil, err
	}


	commentDate := time.Now().UTC().Format(time.RFC3339)

	err := db.insertComment(thisCommentId, MessageId, CommentBody, commentDate)
	if err != nil {
		return nil, err
	}

	newComment = structs.Comment {
		CommentId:            thisCommentId,
		MessageId:            MessageId,
		CommentBody:          CommentBody,
		Date:                 commentDate,
	}

	return newComment, err


	
}

func (db * appdbimpl) uncommentMessage(CommentId structs.Identifier) error {
	_, err:= db.c.Exec(`DELETE FROM comment WHERE CommentId = ?`, CommentId)
	return err
}