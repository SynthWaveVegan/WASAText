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
	"log"
	// "time"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	//TODO
	DoLogin(Username string) (structs.Identifier, error)
	//checkValidId(checkingId string, startId string) (bool, error)
	//checkUserExist(Username string) (string, bool, error)

	getConversation(ConversationId structs.Identifier) (structs.Conversation, error)
	getMyConversations(UserHosting structs.Identifier) ([]structs.Identifier, error)

	SetMyUsername(mode string, newName string, UserId string) error
	//CreateUser(Username string, UserId string) error
	//checkUserExist(Username string) (string, bool, error)

	//sendPhoto(file []byte, format string, UploaderId structs.Identifier) (structs.Photo, error)
	//forwardPhoto(OldPhotoId structs.Identifier, UploaderId structs.Identifier) (structs.Photo, error)
	//insertPhoto(PhotoId string, Path string, UploaderId string, Date string) (error)

	//insertMessage(MessageBody string, Comments []structs.Comment, UploaderUserid structs.Identifier, MessageId structs.Identifier, Date string) (error)
	sendMessage(MessageBody string, UploaderId structs.Identifier) (structs.Message, error)
	forwardMessage(OldMessageId structs.Identifier, UploaderId structs.Identifier) (structs.Message, error)

	commentMessage(CommentBody string, MessageId structs.Identifier, UploaderId structs.Identifier) (structs.Comment, error)
	//insertComment(CommentId structs.Identifier, MessageId structs.Identifier, CommentBody string, Date string) error
	uncommentMessage(CommentId structs.Identifier) error
	deleteMessage(MessageId structs.Identifier) error

	addToGroup(GroupId structs.Identifier, AddUserId structs.Identifier) error
	leaveGroup(GroupId structs.Identifier, UserId structs.Identifier) error
	SetGroupName(mode string, newName string, GroupId structs.Identifier) (string, error)
	//createGroup(GroupId structs.Identifier, User []structs.User, Messages []structs.Message, GroupName string, Photo structs.Photo) error

	SetMyPhoto(userId structs.Identifier, photoLink string) error
	setGroupPhoto(photoLink string, GroupId structs.Identifier) error

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

		userQuery := `CREATE TABLE IF NOT EXISTS user (
		          UserId VARCHAR(11) NOT NULL PRIMARY KEY,
		          Username VARCHAR(16) NOT NULL,
				  UserPhoto TEXT
	)`
		messageQuery := `CREATE TABLE IF NOT EXISTS message (
		MessageId VARCHAR(11) NOT NULL PRIMARY KEY,
		MessageBody TEXT,
		Date TEXT,
		UploaderId VARCHAR(11) NOT NULL,
		FOREIGN KEY (UploaderId) REFERENCES user(UserId)
	)`
		commentQuery := `CREATE TABLE IF NOT EXISTS comment (
		CommentId VARCHAR(11) NOT NULL PRIMARY KEY,
		CommentBody TEXT,
		Date TEXT,
		UploaderId VARCHAR(11) NOT NULL,
		MessageId VARCHAR(11) NOT NULL,
		FOREIGN KEY (UploaderId) REFERENCES user(UserId),
		FOREIGN KEY (MessageId) REFERENCES message(MessageId)
	)`
		userGroupQuery := `CREATE TABLE IF NOT EXISTS userGroup (
		GroupId VARCHAR(11) NOT NULL,
		UserId VARCHAR(11) NOT NULL,
		PRIMARY KEY (GroupId, UserId),
		FOREIGN KEY (UserId) REFERENCES user(UserId),
		FOREIGN KEY (GroupId) REFERENCES groups(GroupId)
		
	)`
		conversationQuery := `CREATE TABLE IF NOT EXISTS conversation (
		ConversationId VARCHAR(11) NOT NULL PRIMARY KEY,
		UserHosting VARCHAR(11) NOT NULL,
		UserConnected VARCHAR(11) NOT NULL,
		FOREIGN KEY (UserConnected) REFERENCES user(UserId),
		FOREIGN KEY (UserHosting) REFERENCES user(UserId)

	)`
		photoQuery := `CREATE TABLE IF NOT EXISTS photo (
		PhotoId VARCHAR(11) NOT NULL PRIMARY KEY,
		UploaderId VARCHAR(11) NOT NULL,
		Date TEXT,
		PhotoPath TEXT,
		FOREIGN KEY (UploaderId) REFERENCES user(UserId)
	)`

		groupQuery := `CREATE TABLE IF NOT EXISTS groups (
		GroupId VARCHAR(11) NOT NULL PRIMARY KEY,
		GroupName VARCHAR(16) NOT NULL,
		GroupPhoto TEXT
	)`

		err = execQueries(db, userQuery, messageQuery, commentQuery, userGroupQuery, conversationQuery, photoQuery, groupQuery)
		if err != nil {
			log.Println("Error creating tables: ", err)
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
