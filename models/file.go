package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"gloom/storage"
	"io/ioutil"
	"os"
	"time"
)

type File struct {
	Id         int
	UploaderId string
	Filename   string
	InsertDate time.Time
}

var filesPath = "static/files/"

var filesMaxGlobal = 1000
var filesMaxUploader = 15

// 5Gb / MaxFiles
var filesMaxSize = int64((1000 * 1000 * 1000 * 5) / filesMaxGlobal)

// 3 Days
var filesExpiration = int64(1 * 60 * 60 * 24 * 3)

func deleteFile(db *sql.DB, file *File) {
	path := filesPath + file.Filename
	err := os.Remove(path)

	if err != nil {
		query := `
		DELETE FROM file p
		WHERE p.id = ?
	`
		storage.PreparedExec(db, query, file.Id)

		panic(err)
	}

	query := `
		DELETE FROM file p
		WHERE p.id = ?
	`
	storage.PreparedExec(db, query, file.Id)
}

func pruneFiles(db *sql.DB, files []File) {
	for _, file := range files {
		age := time.Now().Unix() - file.InsertDate.Unix()
		if age > filesExpiration {
			deleteFile(db, &file)
		}
	}
}

func checkLimits(content []byte, uploaderId string) bool {
	// Length check in bytes is intentional
	if int64(len(content)) > filesMaxSize {
		return false
	}

	files := GetFiles()

	if len(files) > filesMaxGlobal {
		return false
	}

	userFiles := 0
	for _, file := range files {
		if file.UploaderId == uploaderId {
			userFiles++
		}
	}

	if userFiles > filesMaxUploader {
		return false
	}

	return true
}

func GetFiles() []File {
	var files []File

	rows := storage.PreparedQuery(
		storage.Db,
		"SELECT id, uploader_id, filename, insert_date FROM file",
	)
	defer rows.Close()

	for rows.Next() {
		var file File
		var timestamp int64
		err := rows.Scan(
			&file.Id, &file.UploaderId, &file.Filename, &timestamp,
		)
		if err != nil {
			panic(err)
		}

		file.InsertDate = time.Unix(timestamp, 0)
		files = append(files, file)
	}

	return files
}

func GetFile(db *sql.DB, fileId int64) *File {
	query := `
		SELECT
			id,
			uploader_id,
			filename,
			insert_date
		FROM
			file
		WHERE
			id = ?
	`
	row := storage.PreparedQueryRow(db, query, fileId)

	var file File
	row.Scan(&file.Id, &file.UploaderId, &file.Filename, &file.InsertDate)

	return &file
}

func SearchFile(db *sql.DB, filename string) *File {
	query := `
		SELECT
			id,
			uploader_id,
			filename,
			insert_date
		FROM
			file
		WHERE
			filename = ?
	`
	row := storage.PreparedQueryRow(db, query, filename)

	var file File
	row.Scan(&file.Id, &file.UploaderId, &file.Filename, &file.InsertDate)

	return &file
}

func NewFile(db *sql.DB, content []byte, uploaderId string) (*File, error) {
	// Move these outside of the function, into the handler
	pruneFiles(db, GetFiles())

	if !checkLimits(content, uploaderId) {
		return nil, errors.New("file limit reached")
	}

	h := sha256.New()
	h.Write(content)
	filename := hex.EncodeToString(h.Sum(nil))

	query := `
		INSERT INTO file (
			uploader_id,
			filename,
			insert_date
		) VALUES (
			?,
			?,
			?
		)
	`

	fileId, err := storage.PreparedExec(
		db, query, uploaderId, filename, time.Now().Unix(),
	)

	// Hash collided with another file, return the existing one
	if err != nil {
		return SearchFile(db, filename), nil
	}

	err = ioutil.WriteFile(filesPath+filename, content, 0644)
	if err != nil {
		panic(err)
	}

	return GetFile(db, fileId), nil
}
