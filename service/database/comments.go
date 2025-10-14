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

func (db * appdbimpl) insertComment(CommentId structs.Identifier, MessageId structs.Identifier, CommentBody string, Date string, UploaderId structs.Identifier) error {

	_, err := db.c.Exec(`INSERT INTO comment(CommentId, MessageId, CommentBody, Date, UploaderId) VALUES (?, ?, ?, ?, ?)` (CommentId, MessageId, CommentBody, Date, UploaderId))
	
	return err
}

func (db * appdbimpl) commentMessage(CommentBody string, MessageId structs.Identifier, UploaderId structs.Identifier) (structs.Comment, error) {

	var thisCommentId string
	var checkId false

	thisCommentId, err = generateIdentifier("C")
	if err != nil {
		return nil, err
	}
	checkId, err = db.checkValidId(thisCommentId, "C")
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