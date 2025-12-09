package database

import (
	"database/sql"
	"errors"
	// "fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
	"context"
)

func (db *appdbimpl) insertGroup(GroupId structs.Identifier, GroupName string, PhotoPath string) error {
	_, err := db.c.ExecContext(context.Background(),`INSERT INTO groups (GroupId, GroupName, GroupPhoto) VALUES (?, ?, ?)`, GroupId, GroupName, PhotoPath)
	return err
}

func (db *appdbimpl) createGroup(GroupName string, UserId structs.Identifier) (structs.Group, error) {

	var Users []structs.User
	var Messages []structs.Message

	thisGroupId  := generateIdentifier("G")

	validId, err := db.checkValidId(thisGroupId.Id, "G")
	if err != nil {
		return structs.Group{}, err
	}
	if !validId {
		return structs.Group{}, err
	}

	CreatorUser, err := db.GetUser(UserId)
	if err != nil {
		return structs.Group{}, err
	}
	Users = append(Users, CreatorUser)

	err = db.SetGroupName("New", GroupName, thisGroupId)
	if err != nil {
		return structs.Group{}, err
	}

	err = db.insertGroup(thisGroupId, GroupName, "")
	if err != nil {
		return structs.Group{}, err
	}

	newGroup := structs.Group{

		GroupId:   thisGroupId,
		GroupName: GroupName,
		Users:     Users,
		Messages:  Messages,
	}

	return newGroup, nil

}

func (db *appdbimpl) SetGroupName(mode string, newName string, GroupId structs.Identifier) error {

	var counter int
	validName := checkValidName(newName)

	err := db.c.QueryRowContext(context.Background(),`SELECT COUNT(*) FROM groups WHERE GroupName = ?`, newName).Scan(&counter)
	if err != nil {
		return err
	}
	if counter != 0 {
		log.Printf("name taken")
		return nil
	}
	if !validName {
		log.Printf("name invalid")
		return nil
	}
	if counter == 0 && validName {

		switch mode {

		case "New":

			err = db.insertGroup(GroupId, newName, "")
			if err != nil {
				return err
			}
			return nil

		case "Update":

			_, err = db.c.ExecContext(context.Background(),`UPDATE groups SET GroupName = ? WHERE GroupId = ?`, newName, GroupId)
			if err != nil {
				return err
			}
			return nil

		default:
			return err
		}
	}
	return err

}

func (db *appdbimpl) AddToGroup(GroupName string, AddUserId structs.Identifier) error {

	var counter int
	var checkId bool
	var GroupId string

	

	GroupId, checkId, err := db.CheckGroupExist(GroupName)
	if err != nil {
		return err
	}
	if !checkId {
		NewGroup, err := db.createGroup(GroupName, AddUserId)
		if err != nil {
			return err
		}

		err = db.insertUserinGroup(NewGroup.GroupId, AddUserId)
		if err != nil {
			return err
		}

		return nil
	}

	err = db.c.QueryRowContext(context.Background(),`SELECT COUNT(*) FROM userGroup WHERE GroupId = ? AND UserId = ?`, GroupId, AddUserId).Scan(&counter)

	if err != nil {
		return err
	}

	if counter == 1 { //user already in group
		return err
	}

	if counter == 0 {

		thisGroupId := structs.Identifier{
			Id: GroupId,
		}

		err = db.insertUserinGroup(thisGroupId, AddUserId)
		if err != nil {
			return err
		}

		return nil
	}

	return err
}

func (db *appdbimpl) insertUserinGroup(GroupId structs.Identifier, AddUserId structs.Identifier) error {

	_, err := db.c.ExecContext(context.Background(),`INSERT INTO userGroup (GroupId, UserId) VALUES (?, ?)`, GroupId, AddUserId)

	return err
}

func (db *appdbimpl) LeaveGroup(GroupId structs.Identifier, UserId structs.Identifier) error {

	var counter int

	err := db.c.QueryRowContext(context.Background(),`SELECT COUNT(*) FROM userGroup WHERE GroupId = ? AND UserId = ?`, GroupId, UserId).Scan(&counter)

	if err != nil {
		return err
	}
	if counter == 0 {
		return err // user not in group
	}
	if counter == 1 {

		err := db.removeUserFromGroup(GroupId, UserId)
		if err != nil {
			return err
		}

		return nil
	}

	return err

}

func (db *appdbimpl) removeUserFromGroup(GroupId structs.Identifier, UserId structs.Identifier) error {
	_, err := db.c.ExecContext(context.Background(),`DELETE FROM userGroup WHERE GroupId = ? AND UserId = ?`, GroupId, UserId)

	return err
}

func (db *appdbimpl) SetGroupPhoto(photoLink string, GroupId structs.Identifier) error {
	_, err := db.c.ExecContext(context.Background(),`UPDATE groups SET GroupPhoto = ? WHERE Groupid = ?`, photoLink, GroupId)

	return err
}

func (db *appdbimpl) CheckGroupExist(Groupname string) (string, bool, error) {

	var GroupId string

	err := db.c.QueryRowContext(context.Background(),`SELECT GroupId FROM groups WHERE GroupName = ?`, Groupname).Scan(&GroupId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, nil
		} else {
			return "", false, err
		}
	} else {
		return GroupId, true, nil
	}

}
