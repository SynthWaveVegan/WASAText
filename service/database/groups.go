package database

import (
	// "database/sql"
	// "errors"
	// "fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
	// "tim"
)

func (db *appdbimpl) insertGroup(GroupId structs.Identifier, GroupName string) error {
	_, err := db.c.Exec(`INSERT INTO groups (GroupId, GroupName) VALUES (?, ?)`, GroupId, GroupName)
	return err
}

func (db *appdbimpl) createGroup(GroupName string, UserId structs.Identifier) (structs.Group, error) {

	var Users []structs.User

	thisGroupId, err := generateIdentifier("G")
	if err != nil {
		return structs.Group{}, err
	}

	validId, err := db.checkValidId(thisGroupId.Id, "G")
	if err != nil {
		return structs.Group{}, err
	}
	if validId == false {
		return structs.Group{}, err
	}

	CreatorUser, err := db.GetUser(UserId)
	if err != nil {
		return structs.Group{}, err
	}
	Users = append(Users, CreatorUser)

	thisGroupName, err := db.SetGroupName("New", GroupName, thisGroupId)

	if err != nil {
		return structs.Group{}, err
	}

	err = db.insertGroup(thisGroupId, thisGroupName)

	newGroup := structs.Group{

		GroupId:   thisGroupId,
		GroupName: thisGroupName,
		Users:     Users,
	}

	return newGroup, nil

}

func (db *appdbimpl) SetGroupName(mode string, newName string, GroupId structs.Identifier) (string, error) {

	var counter int
	validName := checkValidName(newName)

	err := db.c.QueryRow(`SELECT COUNT(*) FROM groups WHERE GroupName = ?`, newName).Scan(&counter)
	if err != nil {
		return "", err
	}
	if counter != 0 {
		log.Printf("name taken")
		return "", nil
	}
	if validName == false {
		log.Printf("name invalid")
		return "", nil
	}
	if counter == 0 && validName == true {

		switch mode {

		case "New":

			err := db.insertGroup(GroupId, newName)
			return "", err

		case "Update":

			_, err := db.c.Exec(`UPDATE groups SET GroupName = ? WHERE GroupId = ?`, newName, GroupId)
			return newName, err

		default:
			return "", err
		}
	}
	return "", err

}

func (db *appdbimpl) addToGroup(GroupId structs.Identifier, AddUserId structs.Identifier) error {

	var counter int

	err := db.c.QueryRow(`SELECT COUNT(*) FROM userGroup WHERE GroupId = ? AND UserId = ?`, GroupId, AddUserId).Scan(&counter)

	if err != nil {
		return err
	}

	if counter == 1 { //user already in group
		return err
	}

	if counter == 0 {
		err := db.insertUserinGroup(GroupId, AddUserId)
		if err != nil {
			return err
		}

		return nil
	}

	return err
}

func (db *appdbimpl) insertUserinGroup(GroupId structs.Identifier, AddUserId structs.Identifier) error {

	_, err := db.c.Exec(`INSERT INTO userGroup (GroupId, UserId) VALUES (?, ?)`, GroupId, AddUserId)

	return err
}

func (db *appdbimpl) leaveGroup(GroupId structs.Identifier, UserId structs.Identifier) error {

	var counter int

	err := db.c.QueryRow(`SELECT COUNT(*) FROM userGroup WHERE GroupId = ? AND UserId = ?`, GroupId, UserId).Scan(&counter)

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
	_, err := db.c.Exec(`DELETE FROM userGroup WHERE GroupId = ? AND UserId = ?`, GroupId, UserId)

	return err
}
