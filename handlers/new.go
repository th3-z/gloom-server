package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"gloom/models"
	"gloom/storage"

	"github.com/labstack/echo"
)

func New(c echo.Context) error {
	fileData := []byte(c.FormValue("file"))

	h := sha256.New()
	h.Write([]byte(c.RealIP()))
	uploaderId := hex.EncodeToString(h.Sum(nil))

	paste, err := models.NewPaste(storage.Db, fileData, uploaderId)
	if err != nil {
		return err
	}

	return c.Redirect(302, "files/"+paste.Filename)
}
