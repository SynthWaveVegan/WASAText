package database

import (
	//"database/sql"
	//"errors"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"log"
	"context"
)

func (db *appdbimpl) ConversationExists(user1Id, user2Id string) (bool, error) {
	// Assicurati che la query restituisca 'true' solo se entrambi gli utenti sono già presenti nella conversazione.
	query := `
		SELECT COUNT(*) 
		FROM userChat uc1
		JOIN userChat uc2 ON uc1.ConversationId = uc2.ConversationId
		JOIN conversation c ON c.ConversationId = uc1.ConversationId
		WHERE uc1.UserId = ? AND uc2.UserId = ? AND c.IsGroup = FALSE
	`

	var count int
	err := db.c.QueryRow(query, user1Id, user2Id).Scan(&count)
	if err != nil {
		return false, err
	}
	
	if count == 0 {
		return false, nil
	}
	return true, nil
	
}
 
func (db *appdbimpl) CreateConversation(UserHosting structs.Identifier, UserConnectedName string) (structs.Conversation, error) {

	var Messages []structs.Message
	var IsGroup = 0
	

	ConversationId := generateIdentifier("C")
	log.Println("Generated Conversation ID:", ConversationId.Id)
	
	validId, err := db.checkValidId(ConversationId.Id, "C")
	if err != nil {
		return structs.Conversation{}, err
	}
	if !validId {
		return structs.Conversation{}, err
	}

	UserConnected, err := db.GetUserIdByName(UserConnectedName)
	if err != nil {
		return structs.Conversation{}, err
	}

	exists, err := db.ConversationExists(UserHosting.Id, UserConnected.Id)
	log.Println("exist = ", exists)
		if err != nil {
			return structs.Conversation{}, err
		}

		if exists  {
			return structs.Conversation{}, fmt.Errorf("conversation already exists between %s and %s", UserHosting.Id, UserConnected.Id)
		}

	ChatPhoto, err := db.getUserPhotobyId(UserConnected.Id)
	if err != nil {
		return structs.Conversation{}, err
	}
	
	_, err = db.c.ExecContext(context.Background(),`INSERT INTO conversation (ConversationId, ChatPhoto, IsGroup) VALUES (?, ?, ?)`, 
	ConversationId.Id, ChatPhoto, IsGroup)
	if err != nil {
    log.Println("Errore nell'inserimento della conversazione:", err)
    return structs.Conversation{}, err
	}

	_, err = db.c.ExecContext(context.Background(),`INSERT INTO UserChat (ConversationId, UserId) VALUES (?, ?)`, 
	ConversationId.Id, UserConnected.Id)
	if err != nil {
    log.Println("Errore nell'inserimento della conversazione:", err)
    return structs.Conversation{}, err
	}

	_, err = db.c.ExecContext(context.Background(),`INSERT INTO UserChat (ConversationId, UserId) VALUES (?, ?)`, 
	ConversationId.Id, UserHosting.Id)
	if err != nil {
    log.Println("Errore nell'inserimento della conversazione:", err)
    return structs.Conversation{}, err
	}

	var Users []structs.User

	User1, err := db.GetUser(UserHosting)
	if err != nil {
		return structs.Conversation{}, err
	}

	Users = append(Users, User1)

	User2, err := db.GetUser(UserConnected)
	if err != nil {
		return structs.Conversation{}, err
	}

	Users = append(Users, User2)


	NewChat := structs.Conversation{

		ConversationId: ConversationId,
		Users:          Users,
		ChatName:       UserConnectedName,
		Messages:       Messages,
	}

	return NewChat, nil
}

func (db *appdbimpl) GetConversation(ConversationId structs.Identifier, currentUserId string) (structs.Conversation, error) {

	// Recupera tutti gli utenti associati alla conversazione
	rows, err := db.c.QueryContext(context.Background(), `
		SELECT u.UserId, u.Username, u.UserPhoto
		FROM users u
		JOIN userChat uc ON u.UserId = uc.UserId
		WHERE uc.ConversationId = ?`, ConversationId.Id)

	if err != nil {
		log.Println("ERROR QUERY (getting users):", err)
		return structs.Conversation{}, err
	}
	defer rows.Close()

	// Crea una slice di utenti
	var users []structs.User
	var otherUserId string

	for rows.Next() {
		var u structs.User
		err := rows.Scan(&u.UserId.Id, &u.Username, &u.UserPhoto)
		if err != nil {
			log.Println("ERROR SCAN (users):", err)
			return structs.Conversation{}, err
		}
		users = append(users, u)

		// Trova l'altro utente (quello che non è l'attuale)
		if u.UserId.Id != currentUserId {
			otherUserId = u.UserId.Id
		}
	}

	if len(users) == 0 {
		return structs.Conversation{}, fmt.Errorf("no users found in conversation")
	}

	// Recupera il nome dell'altro utente (quello che non sta usando la funzione)
	var chatName string
	if otherUserId != "" {
		// Ottieni il nome dell'altro utente
		OtherUserId := structs.Identifier{
			Id: otherUserId,
		}

		otherUsername, err := db.getUsernamebyId(OtherUserId)
		if err != nil {
			log.Println("ERROR GETTING USER NAME:", err)
			return structs.Conversation{}, err
		}
		chatName = otherUsername
	}

	// Recupera i messaggi della conversazione
	messageRows, err := db.c.QueryContext(context.Background(), `
		SELECT MessageId, MessageBody, Date, UploaderId, MediaType
		FROM message
		WHERE ConversationId = ?
		ORDER BY Date DESC`, ConversationId.Id)

	if err != nil {
		log.Println("ERROR QUERY (getting messages):", err)
		return structs.Conversation{}, err
	}
	defer messageRows.Close()

	var messages []structs.Message

	// Scansione dei messaggi
	for messageRows.Next() {

		var msg structs.Message

		err := messageRows.Scan(&msg.MessageId, &msg.MessageBody, &msg.Date, &msg.UploaderId, &msg.MediaType)
		if err != nil {
			log.Println("ERROR SCAN (messages):", err)
			return structs.Conversation{}, err
		}

		msg.ConversationId = ConversationId

		messages = append(messages, msg)
	}

	// Crea la struct di Conversation
	ChatRetrieved := structs.Conversation{
		ConversationId: ConversationId,
		Users:          users,           // Inserisci tutti gli utenti
		ChatName:       chatName,        // Nome chat (quello dell'altro utente)
		Messages:       messages,        // Messaggi
	}

	return ChatRetrieved, nil
}


func (db *appdbimpl) GetMyConversations(UserHosting structs.Identifier) ([]structs.Identifier, error) {

	var ChatStream []structs.Identifier
	var ConversationId structs.Identifier

	rows, err := db.c.QueryContext(context.Background(),`SELECT ConversationId FROM conversation WHERE UserHosting = ?`, UserHosting.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&ConversationId.Id)
		if err != nil {
			return nil, err
		}
		ChatStream = append(ChatStream, ConversationId)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return ChatStream, nil

}
