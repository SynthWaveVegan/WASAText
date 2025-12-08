package database

/*

import (
	// "database/sql"
	// "errors"
	// "fmt"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	// "log"
	"os"
	"path/filepath"
	"time"
)

func (db *appdbimpl) insertPhoto(PhotoId structs.Identifier, Path string, UploaderId structs.Identifier, Date string) error {
	_, err := db.c.Exec(`INSERT INTO photo (PhotoId, PhotoPath,  Date, UploaderId) VALUES (?, ?, ?, ?)`, PhotoId, Path, Date, UploaderId)

	return err
}

func (db *appdbimpl) sendPhoto(file []byte, format string, UploaderId structs.Identifier) (structs.Photo, error) {

	const Folder string = "/tmp/WASAText/images/"

	thisPhotoId, err := generateIdentifier("P")
	if err != nil {
		return structs.Photo{}, err
	}

	checkId, err := db.checkValidId(thisPhotoId.Id, "P")

	if checkId == false {
		return structs.Photo{}, err
	}
	if err != nil {
		return structs.Photo{}, err
	}

	photoDate := time.Now().UTC().Format(time.RFC3339)

	PhotoPath := Folder + UploaderId.Id + "/" + thisPhotoId.Id + "." + format

	err = savePhoto(file, PhotoPath)
	if err != nil {
		return structs.Photo{}, err
	}

	err = db.insertPhoto(thisPhotoId, PhotoPath, UploaderId, photoDate)

	newPhoto := structs.Photo{

		PhotoId:    thisPhotoId,
		Path:       UploaderId.Id + "/" + thisPhotoId.Id + "." + format,
		UploaderId: UploaderId,
		Date:       photoDate,
		Comments:   []structs.Comment{},
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

func (db *appdbimpl) forwardPhoto(OldPhotoId structs.Identifier, UploaderId structs.Identifier) (structs.Photo, error) {

	var PhotoPath string

	thisPhotoId, err := generateIdentifier("P")
	if err != nil {
		return structs.Photo{}, err
	}

	checkId, err := db.checkValidId(thisPhotoId.Id, "P")

	if checkId == false {
		return structs.Photo{}, err
	}
	if err != nil {
		return structs.Photo{}, err
	}

	photoDate := time.Now().UTC().Format(time.RFC3339)

	err = db.c.QueryRow(`SELECT PhotoPath FROM photo WHERE PhotoId = ?`, OldPhotoId).Scan(&PhotoPath)

	err = db.insertPhoto(thisPhotoId, PhotoPath, UploaderId, photoDate)

	forwardedPhoto := structs.Photo{

		PhotoId:    thisPhotoId,
		Path:       PhotoPath,
		UploaderId: UploaderId,
		Date:       photoDate,
		Comments:   []structs.Comment{},
	}

	return forwardedPhoto, nil

}

*/
