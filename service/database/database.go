/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	

	//TODO 
	//• doLogin 
	getConversation() (structs.Conversation, error)
	//• getMyConversations
	SetMyUserame(name string) error
	insertMessage(MessageBody string, Comments []structs.Comment, UploaderUserid structs.Identifier, MessageId structs.Identifier, Date string, ConversationId structs.Identifier) (error)
	sendMessage(MessageBody string, UploaderId structs.Identifier, ConversationId structs.Identifier) (structs.Message , error)
	//• forwardMessage
	commentMessage(CommentBody string, MessageId structs.Identifier) (structs.Comment, error)
	insertComment(CommentId structs.Identifier, MessageId structs.Identifier, CommentBody string, Date string) error
	uncommentMessage(CommentId structs.Identifier) error 
	deleteMessage(messageId string) error 
	//• addToGroup
	//• leaveGroup
	//• setGroupName
	createGroup(GroupId structs.Identifier, User []structs.User, Messages []structs.Message, GroupName string, Photo structs.Photo) error
	//• setMyPhoto
	//• setGroupPhoto

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		
		userQuery := `CREATE TABLE IF NOT EXIST user (
		          userId VARCHAR(11) NOT NULL PRIMARY KEY,
		          username VARCHAR(16) NOT NULL
	)`
		messageQuery := `CREATE TABLE IF NOT EXIST message (
		messageId VARCHAR(11) NOT NULL PRIMARY KEY,
		messageBody TEXT,
		date TEXT,
		uploaderId VARCHAR(11) NOT NULL,
		conversationId VARCHAR(11) NOT NULL,
		FOREIGN KEY (uploaderId) REFERENCES user(userId),
		FOREIGN KEY (conversationId) REFERENCES conversation(conversationId)
	)`
		commentQuery := `CREATE TABLE IF NOT EXIST comment (
		commentId VARCHAR(11) NOT NULL PRIMARY KEY,
		commentBody TEXT,
		date TEXT,
		uploaderId VARCHAR(11) NOT NULL,
		messageId VARCHAR(11) NOT NULL,
		FOREIGN KEY (uploaderId) REFERENCES user(userId),
		FOREIGN KEY (messageId) REFERENCES message(messageId)
	)`
		groupQuery := `CREATE TABLE IF NOT EXIST group (
		groupId VARCHAR(11) NOT NULL PRIMARY KEY,
		groupName VARCHAR(16) NOT NULL

	)`
		conversationQuery := `CREATE TABLE IF NOT EXIST conversation (
		conversationId VARCHAR(11) NOT NULL PRIMARY KEY,
		userConnected VARCHAR(11),
		FOREIGN KEY (userConnected) REFERENCES user(userId)

	)`
		//da creare photo table

		err = execQueries(db, userQuery, messageQuery, commentQuery, groupQuery, conversationQuery)
		if err != nil {
			log.Println("Error creating tables")
		}
		
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func execQueries(db *sql.DB, tables ...string) error {
	for _, el := range tables {
		_, err := db.Exec(el)

		if err != nil {

			return fmt.Errorf("error creating db structurev: %w", err)
		}
	}
	return nil
}
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
