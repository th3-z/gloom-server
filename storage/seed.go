package storage

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"
)

const userQuery = `
	INSERT INTO user 
		(name, password, admin, insert_date)
	VALUES
		(?, ?, ?, ?)
`

func SeedDb(db *sql.DB) {
	h := sha256.New()
	h.Write([]byte("admin" + Salt))
	adminPassword := hex.EncodeToString(h.Sum(nil))

	_, err := PreparedExec(
		db, userQuery,
		"admin", adminPassword, 1, time.Now().Unix(),
	)

	if err != nil {
		panic(err)
	}
}
