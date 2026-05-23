 package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
	"context"
)

func (db *appdbimpl) insertGroup(GroupId structs.Identifier, GroupName string, PhotoPath string) error {
	var IsGroup = "Yes"
	_, err := db.c.ExecContext(context.Background(),`INSERT INTO conversation (ConversationId, ChatPhoto, GroupName, IsGroup) VALUES (?, ?, ?, ?)`,
	 GroupId.Id, PhotoPath, GroupName, IsGroup)
	return err
}

func (db *appdbimpl) CreateGroup(GroupName string, CreatorUserId structs.Identifier, ConversationId string) (structs.Group, error) {

	var Users []structs.User
	var Messages []structs.Message

	thisGroupId  := generateIdentifier("C")
	log.Printf("groupId: %s", thisGroupId)

	validId, err := db.checkValidId(thisGroupId.Id, "C")
	if err != nil {
		return structs.Group{}, err
	}
	if !validId {
		return structs.Group{}, err
	}

	CreatorUser, err := db.GetUser(CreatorUserId)
	if err != nil {
		return structs.Group{}, err
	}
	Users = append(Users, CreatorUser)

	rows, err := db.c.QueryContext(context.Background(), `
		SELECT UserId FROM userChat WHERE ConversationId = ?`, ConversationId)

	if err != nil {
		log.Println("ERROR QUERY (getting users):", err)
		return structs.Group{}, err
	}
	defer rows.Close()

	var AddUserId string

	for rows.Next() {
		var thisUserId string
		err := rows.Scan(&thisUserId)
		if err != nil {
			log.Println("ERROR SCAN (users):", err)
			return structs.Group{}, err
		}

		if thisUserId != CreatorUserId.Id {
			AddUserId = thisUserId
		}
	}

	AddUser := structs.Identifier{
		Id: AddUserId,
	}

	User2, err := db.GetUser(AddUser)
	if err != nil {
		return structs.Group{}, err
	}

	Users = append(Users, User2)

	err = db.SetGroupName("New", GroupName, thisGroupId)
	log.Printf("set group name: %s", GroupName)
	if err != nil {
		return structs.Group{}, err
	}

	for _, user := range Users {
		err = db.insertUserinGroup(thisGroupId, user.UserId)
	}
	
	var GroupPhoto = ""

	newGroup := structs.Group{

		GroupId:   thisGroupId,
		GroupName: GroupName,
		Users:     Users,
		Messages:  Messages,
		ChatPhoto: GroupPhoto,
	}

	return newGroup, nil

}

func (db *appdbimpl) SetGroupName(mode string, newName string, GroupId structs.Identifier) error {

	var counter int
	validName := checkValidName(newName)

	err := db.c.QueryRowContext(context.Background(),`SELECT COUNT(*) FROM conversation WHERE GroupName = ?`, newName).Scan(&counter)
	if err != nil {
		return err
	}
	if counter != 0 {
		return fmt.Errorf("group name already taken")
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

			_, err = db.c.ExecContext(context.Background(),`UPDATE conversation SET GroupName = ? WHERE ConversationId = ?`, newName, GroupId.Id)
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
/*func (db *appdbimpl) AddToGroup(GroupName string, CreatorUser structs.Identifier) error {

	var counter int
	var checkId bool
	var GroupId string

	// Verifica se il gruppo esiste già
	GroupId, checkId, err := db.CheckGroupExist(GroupName)
	if err != nil {
		return fmt.Errorf("error checking if group exists: %w", err)
	}

	if !checkId {
		// Se il gruppo non esiste, crealo
		NewGroup, err := db.CreateGroup(GroupName, CreatorUser)
		log.Printf("Returning group: %#v", NewGroup)
		if err != nil {
			return fmt.Errorf("error creating new group: %w", err)
		}

		// Aggiungi l'utente al gruppo
		err = db.insertUserinGroup(NewGroup.GroupId, CreatorUser)
		log.Printf("2 groupId: %s", NewGroup.GroupId)
		if err != nil {
			return fmt.Errorf("error inserting user into new group: %w", err)
		}
		GroupId = NewGroup.GroupId.Id

		
	}
	rows, err := db.c.QueryContext(context.Background(), `
		SELECT UserId FROM userChat WHERE ConversationId = ?`, ConversationId)

	if err != nil {
		log.Println("ERROR QUERY (getting users):", err)
		return err
	}
	defer rows.Close()

	//var users []structs.User
	var AddUserId string

	for rows.Next() {
		var thisUserId string
		err := rows.Scan(&thisUserId)
		if err != nil {
			log.Println("ERROR SCAN (users):", err)
			return err
		}
		//users = append(users, u)

		if thisUserId != CreatorUser.Id {
			AddUserId = thisUserId
		}
	}

	AddUser := structs.Identifier{
		Id: AddUserId,
	}
	
	// Se il gruppo esiste, aggiungi l'utente
	err = db.c.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM userChat WHERE ConversationId = ? AND UserId = ?`, GroupId, AddUser.Id).Scan(&counter)
	if err != nil {
		return fmt.Errorf("error checking user membership in group: %w", err)
	}

	if counter == 1 {
		// Se l'utente è già nel gruppo
		return fmt.Errorf("user is already a member of the group")
	}

	if counter == 0 {
		thisGroupId := structs.Identifier{Id: GroupId}
		log.Printf("Adding user %s to group %s", AddUser.Id, GroupId)
		err = db.insertUserinGroup(thisGroupId, AddUser)
		if err != nil {
			return fmt.Errorf("error inserting user into existing group: %w", err)
		}
		return nil
	}

	return fmt.Errorf("unexpected error adding user to group")
}*/


func (db *appdbimpl) insertUserinGroup(GroupId structs.Identifier, AddUserId structs.Identifier) error {

	_, err := db.c.ExecContext(context.Background(),`INSERT INTO userChat (ConversationId, UserId) VALUES (?, ?)`, GroupId.Id, AddUserId.Id)

	return err
}

func (db *appdbimpl) LeaveGroup(GroupId structs.Identifier, UserId structs.Identifier) error {

	var counter int

	err := db.c.QueryRowContext(context.Background(),`SELECT COUNT(*) FROM userChat WHERE ConversationId = ? AND UserId = ?`, GroupId.Id, UserId.Id).Scan(&counter)

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
	_, err := db.c.ExecContext(context.Background(),`DELETE FROM userChat WHERE ConversationId = ? AND UserId = ?`, GroupId.Id, UserId.Id)

	return err
}

func (db *appdbimpl) SetGroupPhoto(photoLink string, GroupId structs.Identifier) error {
	_, err := db.c.ExecContext(context.Background(),`UPDATE conversation SET ChatPhoto = ? WHERE ConversationId = ?`, photoLink, GroupId.Id)

	return err
}

func (db *appdbimpl) CheckGroupExist(Groupname string) (string, bool, error) {

	var GroupId string

	err := db.c.QueryRowContext(context.Background(),`SELECT ConversationId FROM conversation WHERE GroupName = ?`, Groupname).Scan(&GroupId)

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
