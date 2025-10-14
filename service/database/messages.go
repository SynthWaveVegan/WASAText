package database

import (
	"errors"
	"database/sql"
	"log"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)

func (db * appdbimpl) insertMessage(MessageBody string, UploaderId structs.Identifier, MessageId structs.Identifier, Date string) (error) {
	_, err := db.c.Exec(`INSERT INTO message (MessageId, MessageBody,  Date, UploaderId) VALUES (?, ?, ?, ?)` (MessageId, MessageBody,  Date,  UploaderId))

	return err
}

func (db * appdbimpl) sendMessage(MessageBody string, UploaderId structs.Identifier) (structs.Message , error) {

	var thisMessageId structs.Identifier
	var checkId false

	thisMessageId, err = generateIdentifier("M")
	if err != nil {
		return nil, err
	}
	
	checkId, err = db.checkValidId(thisMessageId, "M")

	if checkId == false {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	messageDate := time.Now().UTC().Format(time.RFC3339)

	err := db.insertMessage(MessageBody, UploaderId, thisMessageId, messageDate)
	if err != nil {
		return nil, err
	}


	newMessage = structs.Message {
		MessageBody:          MessageBody,
		Comments:             []structs.Comment{},
		UploaderId:           UploaderId,
		MessageId:            thisMessageId,
		Date:                 messageDate,
	}

	return newMessage, nil
}

func (db * appdbimpl) forwardMessage(OldMessageId structs.Identifier, UploaderId structs.Identifier) (structs.Message, error) {

	startId := first_character(OldMessageId)
	if startId =="M"{

		var MessageBody string

		thisMessageId, err = generateIdentifier("M")
		if err != nil {
			return nil, err
		}
	
		checkId, err = db.checkValidId(thisMessageId, "M")

		if checkId == false {
			return nil, err
		}
		if err != nil {
			return nil, err
		}

		messageDate := time.Now().UTC().Format(time.RFC3339)
	
		err := db.c.QueryRow(`SELECT MessageBody FROM message WHERE MessageId = ?`, OldMessageId).Scan(&MessageBody)


		err := db.insertMessage(MessageBody, []structs.Comment{}, UploaderId, thisMessageId, messageDate)
		if err != nil {
			return nil, err
		}


		forwardedMessage = structs.Message {
			MessageBody:          MessageBody,
			Comments:             []structs.Comment{},
			UploaderId:           UploaderId,
			MessageId:            thisMessageId,
			Date:                 messageDate,
		}	

		return forwardedMessage, nil

	}else if startId == "P" { //se il messaggio è una foto

		var PhotoPath string

		thisPhotoId, err = generateIdentifier("P")
		if err != nil {
			return nil, err
		}
	
		checkId, err = db.checkValidId(thisPhotoId, "P")

		if checkId == false {
			return nil, err
		}
		if err != nil {
			return nil, err
		}

		photoDate := time.Now().UTC().Format(time.RFC3339)

		err := db.c.QueryRow(`SELECT PhotoPath FROM photo WHERE PhotoId = ?`, OldMessageId).Scan(&PhotoPath)

		err := db.insertPhoto(thisPhotoId, PhotoPath, UploaderId, photoDate)

		forwardedPhoto = structs.Photo {

			PhotoId:              thisPhotoId,
			Path:                 UploaderId + "/" + thisPhotoId "." + format,
			UploaderId:           UploaderId,
			Date:                 photoDate,
			Comments:             []structs.Comment{},

		}
	
		return forwardedPhoto, nil


	}

}


func (db * appdbimpl) deleteMessage(MessageId structs.Identifier) error {
	_, err := db.c.Exec(`DELETE FROM message WHERE MessageId = ?`, MessageId)
	return err
}

