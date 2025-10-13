package database

import (
	"errors"
	"database/sql"
	"log"
	"time"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)

func (db * appdbimpl) insertGroup(GroupId structs.Identifier, User []structs.User, Messages []structs.Message, GroupName string, Photo structs.Photo) error {
	_, err := db.c.Exec(`INSERT INTO group (GroupId, User, Messages, GroupName, Photo) VALUES (?, ?, ?, ?, ?)` (GroupId, User, Messages, GroupName, Photo))
	return err
}

func (db * appdbimpl) createGroup(GroupName string, UserId structs.Identifier) (structs.Group, error) {

	thisGroupId, error := generateIdentifier("G")
	if err != nil {
		return structs.Group{}, err
	}

	validId, error := checkValidId(thisGroupId, "G")
	if err != nil {
		return structs.Group{}, err
	}
	if validId == false {
		return structs.Group{}, err
	}

	CreatorUser, err := GetUser(UserId)
	if err != nil {
		return structs.Group{}, err
	}
	err := insertGroup(thisGroupId, , []structs.Message{})

type Group struct {

	GroupId              Identifier
	Users                []User
	Messages             []Message
	GroupName            string
	Photo                Photo
}


}
func (db * appdbimpl) getGroupStream()