package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"fmt"
	"gloom/models"
	"gloom/storage"

	"github.com/labstack/echo"
)

// curl -i -X 'POST' -F 'file=@README.md' -F 'password=admin' -F 'username=admin' 'localhost:5001/upload'
func Upload(c echo.Context) error {
	username := c.FormValue("username")
	h := sha256.New()
	h.Write([]byte(c.FormValue("password") + storage.Salt))
	password := hex.EncodeToString(h.Sum(nil))

	user, err := models.LoadUser(storage.Db, username, password)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Bad username or password\n")
	}

	print(user.Name)

	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "No file provided\n")
	}

	models.NewFile(storage.Db, user.Id, file, "/home/the_z/test.md")

	return c.String(http.StatusOK, "File uploaded\n")
}

func List(c echo.Context) error {
	username := c.FormValue("username")
	h := sha256.New()
	h.Write([]byte(c.FormValue("password") + storage.Salt))
	password := hex.EncodeToString(h.Sum(nil))

	user, err := models.LoadUser(storage.Db, username, password)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Bad username or password\n")
	}

	for _, file := range user.Files {
		fmt.Println(file.Path)
	}

	// TODO: Return json of user's files
	return c.String(http.StatusOK, "")
}

func Delete(c echo.Context) error {
	username := c.FormValue("username")
	h := sha256.New()
	h.Write([]byte(c.FormValue("password") + storage.Salt))
	password := hex.EncodeToString(h.Sum(nil))

	user, err := models.LoadUser(storage.Db, username, password)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Bad username or password\n")
	}

	print(user.Name)

	// TODO: Delete requested file
	return c.String(http.StatusOK, "")
}
