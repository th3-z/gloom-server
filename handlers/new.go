package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type TestModel struct {
	file []byte
}

func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

// curl -i -X 'POST' -F 'file=@file.x' 'localhost:5001/new'
func New(c echo.Context) (err error) {
	fileData, _ := c.FormFile("file")

	test := new(TestModel)
	if err = c.Bind(test); err != nil {
		panic(err)
	}

	fmt.Println(CToGoString(test.file[:]))
	fmt.Println(fileData.Filename)

	return c.String(http.StatusOK, "File uploaded\n")
}
