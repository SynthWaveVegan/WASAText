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

func (db *appdbimpl) insertComment(CommentId structs.Identifier, MessageId structs.Identifier, CommentBody string, Date string, UploaderId structs.Identifier) error {

	_, err := db.c.ExecContext(context.Background(),`INSERT INTO comment(CommentId, MessageId, CommentBody, Date, UploaderId) VALUES (?, ?, ?, ?, ?)`, CommentId, MessageId, CommentBody, Date, UploaderId)

	return err
}

func (db *appdbimpl) CommentMessage(CommentBody string, MessageId structs.Identifier, UploaderId structs.Identifier) (structs.Comment, error) {


	thisCommentId := generateIdentifier("C")
	
	checkId, err := db.checkValidId(thisCommentId.Id, "C")
	if !checkId {
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
	_, err := db.c.ExecContext(context.Background(),`DELETE FROM comment WHERE CommentId = ?`, CommentId)
	return err
}
