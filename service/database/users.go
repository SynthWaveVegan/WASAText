package database

import (
	"errors"
	"database/sql"
	"log"
	"fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)

func (db *appdbimpl) DoLogin(username string) error {
	
}
func (db *appdbimpl) SetMyUserame(name string) error {
	_, err := db.c.Exec("INSERT INTO user (username) VALUES (?)", username)
	return err
}

