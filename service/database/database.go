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
	sendMessage(MessageBody string Comments []Comment UploaderUserid Identifier MessageId Identifier Date string ConversationId Identifier) (error)
	//• forwardMessage
	//• commentMessage
	//• uncommentMessage
	deleteMessage(messageId string) error 
	//• addToGroup
	//• leaveGroup
	//• setGroupName
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
		
		createUserTable(db)
		createMessageTable(db)
		createCommentTable(db)
		createGroupTable(db)
		createConversationTable(db)
		
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func createUserTable(db * sql.DB) {
	userQuery := CREATE TABLE IF NOT EXIST user (
		userId VARCHAR(11) NOT NULL PRIMARY KEY,
		username VARCHAR(16) NOT NULL,
	)
	_, err := db.Exec(userQuery)

	if err != nil {
			log.fatal(err)
		}
}

func createMessageTable(db * sql.DB) {
	messageQuery := CREATE TABLE IF NOT EXIST message (
		messageId VARCHAR(11) NOT NULL PRIMARY KEY,
		messageBody TEXT,
		date TEXT,
		uploaderId VARCHAR(11) NOT NULL,
		conversationId VARCHAR(11) NOT NULL,
		FOREIGN KEY (uploaderId) REFERENCES user(userId),
		FOREIGN KEY (conversationId) REFERENCES conversation(conversationId),
	)
	_, err := db.Exec(messageQuery)

	if err != nil {
			log.fatal(err)
		}
}

func createCommentTable(db * sql.DB) {
	commentQuery := CREATE TABLE IF NOT EXIST comment (
		commentId VARCHAR(11) NOT NULL PRIMARY KEY,
		commentBody TEXT,
		date TEXT,
		uploaderId VARCHAR(11) NOT NULL,
		messageId VARCHAR(11) NOT NULL,
		FOREIGN KEY (uploaderId) REFERENCES user(userId),
		FOREIGN KEY (messageId) REFERENCES message(messageId),
	)
	_, err := db.Exec(commentQuery)

	if err != nil {
			log.fatal(err)
		}
}

func createGroupTable(db * sql.DB) {
	groupQuery := CREATE TABLE IF NOT EXIST group (
		groupId VARCHAR(11) NOT NULL PRIMARY KEY,
		groupName VARCHAR(16) NOT NULL,

	)
	_, err := db.Exec(groupQuery)

	if err != nil {
			log.fatal(err)
		}
}

func createConversationTable(db, * sql.DB) {
	conversationQuery := CREATE TABLE IF NOT EXIST conversation (
		conversationId VARCHAR(11) NOT NULL PRIMARY KEY,
		userConnected VARCHAR(11),
		FOREIGN KEY (userConnected) REFERENCES user(userId),

	)
}
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
