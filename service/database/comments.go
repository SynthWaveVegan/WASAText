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

func (db *appdbimpl) insertComment(CommentId structs.Identifier, MessageId structs.Identifier, CommentBody string, Date string, UploaderId structs.Identifier) error {

	_, err := db.c.Exec(`INSERT INTO comment(CommentId, MessageId, CommentBody, Date, UploaderId) VALUES (?, ?, ?, ?, ?)`, CommentId, MessageId, CommentBody, Date, UploaderId)

	return err
}

func (db *appdbimpl) CommentMessage(CommentBody string, MessageId structs.Identifier, UploaderId structs.Identifier) (structs.Comment, error) {

	var thisCommentId structs.Identifier
	var checkId = false

	thisCommentId, err := generateIdentifier("C")
	if err != nil {
		return structs.Comment{}, err
	}
	checkId, err = db.checkValidId(thisCommentId.Id, "C")
	if checkId == false {
		return structs.Comment{}, err
	}
	if err != nil {
		return structs.Comment{}, err
	}

	commentDate := time.Now().UTC().Format(time.RFC3339)

	err = db.insertComment(thisCommentId, MessageId, CommentBody, commentDate, UploaderId)
	if err != nil {
		return structs.Comment{}, err
	}

	newComment := structs.Comment{
		CommentId:   thisCommentId,
		MessageId:   MessageId,
		CommentBody: CommentBody,
		Date:        commentDate,
		UploaderId:  UploaderId,
	}

	return newComment, err

}

func (db *appdbimpl) UncommentMessage(CommentId structs.Identifier) error {
	_, err := db.c.Exec(`DELETE FROM comment WHERE CommentId = ?`, CommentId)
	return err
}
