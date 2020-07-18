package storage

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"
)

const userQuery = `
	INSERT INTO user 
		(name, api_key, insert_date)
	VALUES
		(?, ?, ?)
`

func SeedDb(db *sql.DB) {
	h := sha256.New()
	h.Write([]byte("admin"))
	adminApiKey := hex.EncodeToString(h.Sum(nil))

	_, err := PreparedExec(
		db, userQuery,
		"admin", adminApiKey, time.Now().Unix(),
	)

	if err != nil {
		panic(err)
	}
}
