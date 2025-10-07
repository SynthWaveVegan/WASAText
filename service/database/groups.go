package database

import (
	"errors"
	"database/sql"
	"log"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)

func (db * appdbimpl) createGroup(GroupId structs.Identifier, User []structs.User, Messages []structs.Message, GroupName string, Photo structs.Photo) error {
	_, err := db.c.Exec("INSERT INTO group (GroupId, User, Messages, GroupName, Photo) VALUES (?, ?, ?, ?, ?)" (GroupId, User, Messages, GroupName, Photo))
	return err
}

func (db * appdbimpl) getGroupStream()