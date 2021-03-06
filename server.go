package main

import (
	"fmt"
	"time"

	"github.com/labstack/echo"

	"gloom/config"
	"gloom/handlers"
	"gloom/storage"
)

func routes(e *echo.Echo) {
	e.POST("/upload", handlers.Upload)
	e.POST("/list", handlers.List)

	e.Static("/standing", "static/standing")
	e.Static("/transient", "static/transient")
}

func main() {
	t := time.Now().UTC()

	globalConfig, err := config.FromFile("app.json")
	if err != nil {
		panic(err)
	}
	config.GlobalConfig = globalConfig

	storage.Db = storage.InitDB("storage.db")
	defer storage.Db.Close()
	storage.CreateSchema(storage.Db)
	storage.SeedDb(storage.Db)

	e := echo.New()
	e.Debug = true

	routes(e)

	fmt.Print("Start time: ", t.Format("Mon Jan 2 15:04:05"))
	e.Logger.Fatal(e.Start(":5001"))
}
