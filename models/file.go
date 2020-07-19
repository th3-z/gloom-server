package models

import (
	"database/sql"
	"fmt"
	"gloom/storage"
	"io"
	"mime/multipart"
	"os"
	"time"
)

type File struct {
	Id         int64
	UserId     int64
	Path       string
	InsertDate int64
}

func GetUserFiles(db *sql.DB, userId int64) []*File {
	query := `
		SELECT
			f.id,
			f.user_id,
			f.path,
			f.insert_date
		FROM
			file f
		WHERE
			f.user_id = ?
			AND f.deleted_date IS NULL
	`

	rows := storage.PreparedQuery(
		db, query,
		userId,
	)

	var files []*File
	for rows.Next() {
		var file File
		err := rows.Scan(
			&file.Id, &file.UserId, &file.Path, &file.InsertDate,
		)
		if err != nil {
			fmt.Println(err)
			continue
		}

		files = append(files, &file)
	}

	return files
}

func NewFile(db *sql.DB, userId int64, fileHeader *multipart.FileHeader, dstPath string) (*File, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	insertDate := time.Now().Unix()
	query := `
		INSERT INTO file(
			path,
			user_id,
			insert_date
		) VALUES (?, ?, ?)
	`

	fileId, err := storage.PreparedExec(
		db, query,
		dstPath, userId, insertDate,
	)

	if err != nil {
		return nil, err
	}

	file := File{
		fileId, userId, dstPath, insertDate,
	}

	return &file, nil
}
