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

func (db * appdbimpl) insertPhoto(PhotoId string, Path string, UploaderId string, Date string) (error) {
	_, err := db.c.Exec(`INSERT INTO photo (PhotoId, PhotoPath,  Date, UploaderId) VALUES (?, ?, ?, ?)` (PhotoId, Path,  Date,  UploaderId))

	return err
}

func (db * appdbimpl) createPhoto(file []byte, format string, UploaderId structs.Identifier) (structs.Photo, error) {

	const Folder string = "/tmp/wasatext/WASAText/images/"

	thisPhotoId, err := generateIdentifier("P")
	if err != nil {
		return structs.Photo{}, err
	}

	checkId, err = checkValidId(thisPhotoId, "P")

	if checkId == false {
		return structs.Photo{}, err
	}
	if err != nil {
		return structs.Photo{}, err
	}

	photoDate := time.Now().UTC().Format(time.RFC3339)

	PhotoPath = Folder + UploaderId + "/" + thisPhotoId "." + format

	err := savePhoto(file, PhotoPath)
	if err != nil {
		return structs.Photo{}, err
	}

	newPhoto = structs.Photo {

		PhotoId:              thisPhotoId,
		Path:                 UploaderId + "/" + thisPhotoId "." + format,
		UploaderId:           UploaderId,
		Date:                 photoDate,
		Comments:             []structs.Comment{},

	}
	
	return newPhoto, nil

	
}

func savePhoto(file []byte, path string) error {
	
	dir := filepath.Dir(path)
	
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, file, 0644)

	return err
}


func (db * appdbimpl) sendMessage(MessageBody string, UploaderId structs.Identifier) (structs.Message , error) {

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
	
		checkId, err = checkValidId(thisMessageId, "M")

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
	
		checkId, err = checkValidId(thisPhotoId, "P")

		if checkId == false {
			return nil, err
		}
		if err != nil {
			return nil, err
		}

		photoDate := time.Now().UTC().Format(time.RFC3339)

		err := db.c.QueryRow(`SELECT PhotoPath FROM photo WHERE PhotoId = ?`, OldMessageId).Scan(&PhotoPath)

	}

	

}

//creare tutte le interazioni per messaggi foto




func (db * appdbimpl) deleteMessage(MessageId structs.Identifier) error {
	_, err := db.c.Exec(`DELETE FROM message WHERE MessageId = ?`, MessageId)
	return err
}

func (db * appdbimpl) insertComment(CommentId structs.Identifier, MessageId structs.Identifier, CommentBody string, Date string, UploaderId structs.Identifier) error {

	_, err := db.c.Exec(`INSERT INTO comment(CommentId, MessageId, CommentBody, Date, UploaderId) VALUES (?, ?, ?, ?, ?)` (CommentId, MessageId, CommentBody, Date, UploaderId))
	
	return err
}

func (db * appdbimpl) commentMessage(CommentBody string, MessageId structs.Identifier, UploaderId structs.Identifier) (structs.Comment, error) {

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
		UploaderId:           UploaderId,
	}

	return newComment, err


	
}

func (db * appdbimpl) uncommentMessage(CommentId structs.Identifier) error {
	_, err:= db.c.Exec(`DELETE FROM comment WHERE CommentId = ?`, CommentId)
	return err
}