package database

import (
	"errors"
	"database/sql"
	"log"
	"time"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)

func (db * appdbimpl) insertGroup(GroupId structs.Identifier, GroupName string) error {
	_, err := db.c.Exec(`INSERT INTO group (GroupId, GroupName) VALUES (?, ?)` (GroupId, GroupName))
	return err
}

func (db * appdbimpl) createGroup(GroupName string, UserId structs.Identifier) (structs.Group, error) {

	var Users []structs.User{}

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
	Users = append(Users, CreatorUser)

	thisGroupName, err := db.setGroupName("New", GroupName, thisGroupId)

	if err != nil {
		return structs.Group{}, err
	}

	err := insertGroup(thisGroupId, thisGroupName)

	newGroup: structs.Group {

	GroupId:             thisGroupId,
	GroupName:           thisGroupName,
	Users:				 Users,
	}

	return newGroup, nil

}

func (db *appdbimpl) SetGroupName(mode string, newName string, GroupId string) error {

	var counter int
	validName, err = checkValidName(newName)

	err := db.c.QueryRow(`SELECT COUNT(*) FROM group WHERE GroupName = ?`, newName).Scan(&counter)
	if err != nil {
		return err
	}
	if counter != 0 {
		log.Printf("name taken")
		return nil
	}
	if validName == false {
		log.Printf("name invalid")
		return nil
	}
	if counter == 0 && validName == true {

		switch mode {

		case "New":

			err := db.insertUser(newName, GroupId)
			return err

		case "Update":

			
			_, err := db.c.Exec(`UPDATE group SET GroupName = ? WHERE GroupId = ?`, newName, GroupId)
			return err

		default:
			return err
		}
	}
	
	return nil
}

func (db * appdbimpl) addToGroup(GroupId structs.Identifier, AddUserId structs.Identifier) (error) {

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


}

func (db *appdbimpl) insertUserinGroup(GroupId structs.Identifier, AddUserId structs.Identifier) error {

	_, err := db.c.Exec(`INSERT INTO userGroup (GroupId, UserId) VALUES (?, ?)`, GroupId, AddUserId)

	return err
}

func (db * appdbimpl) leaveGroup(GroupId structs.Identifier, UserId structs.Identifier) error {

	var counter int
	

	err := db.c.QueryRow(`SELECT COUNT(*) Ffunc (db *appdbimpl) SetMyUserame(mode string, newName string, userId string) error {

	var counter int
	validName, err = checkValidName(newName)

	err := db.c.QueryRow(`SELECT COUNT(*) FROM user WHERE username = ?`, newName).Scan(&counter)
	if err != nil {
		return err
	}
	if count != 0 {
		log.Printf("username taken")
		return nil
	}
	if validName == false {
		log.Printf("username invalid")
		return nil
	}
	if count == 0 && validName == true {

		switch mode {

		case "New":

			err := db.createUser(newName, userId)
			return err

		case "Update":

			
			_, err := db.c.Exec(`UPDATE user SET username = ? WHERE userId = ?`, newName, userId)
			return err

		default:
			return err
		}
	}
	
	return nil
}ROM userGroup WHERE GroupId = ? AND UserId = ?`, GroupId, UserId).Scan(&counter)

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


}

func (db * appdbimpl) removeUserFromGroup(GroupId structs.Identifier, UserId structs.Identifier) error {
	_, err := db.c.Exec(`DELETE FROM userGroup WHERE GroupId = ? AND UserId = ?`, GroupId, UserId)

	return err
}