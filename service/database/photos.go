package database

import (
	"errors"
	"database/sql"
	"log"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"github.com/SynthWaveVegan/WASAText/service/structs"
)

func (db * appdbimpl) insertPhoto(PhotoId string, Path string, UploaderId string, Date string) (error) {
	_, err := db.c.Exec(`INSERT INTO photo (PhotoId, PhotoPath,  Date, UploaderId) VALUES (?, ?, ?, ?)` (PhotoId, Path,  Date,  UploaderId))

	return err
}

func (db * appdbimpl) sendPhoto(file []byte, format string, UploaderId structs.Identifier) (structs.Photo, error) {

	const Folder string = "/tmp/wasatext/WASAText/images/"

	thisPhotoId, err := generateIdentifier("P")
	if err != nil {
		return structs.Photo{}, err
	}

	checkId, err = db.checkValidId(thisPhotoId, "P")

	if checkId == false {
		return structs.Photo{}, err
	}
	if err != nil {
		return structs.Photo{}, err
	}

	photoDate := time.Now().UTC().Format(time.RFC3339)

	PhotoPath = Folder + UploaderId + "/" + thisPhotoId "." + format

	err := savePhoto(file, PhotoPath)
	if err != nil {
		return structs.Photo{}, err
	}

	err := db.insertPhoto(thisPhotoId, PhotoPath, UploaderId, photoDate)

	newPhoto = structs.Photo {

		PhotoId:              thisPhotoId,
		Path:                 UploaderId + "/" + thisPhotoId "." + format,
		UploaderId:           UploaderId,
		Date:                 photoDate,
		Comments:             []structs.Comment{},

	}
	
	return newPhoto, nil

	
}

func savePhoto(file []byte, path string) error {
	
	dir := filepath.Dir(path)
	
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, file, 0644)

	return err
}
