package main

import (
	"fmt"

	"github.com/feanor306/image_tagger/src/config"
	"github.com/feanor306/image_tagger/src/database"
	"github.com/feanor306/image_tagger/src/handler"
	"github.com/labstack/echo/v4"
)

// main will start the server
func main() {
	conf := config.GetConfig()

	db, err := database.GetDatabase(conf)
	if err != nil {
		panic(err)
	}

	err = db.InitDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	srv := echo.New()

	tagHandler := handler.TagHandler{
		DB: db,
	}
	mediaHandler := handler.MediaHandler{
		DB: db,
	}

	srv.GET("/", tagHandler.HandleTagShowAll)
	srv.GET("/tags", tagHandler.HandleTagShowAll)
	srv.POST("/tags", tagHandler.HandleTagCreate)
	srv.GET("/media-new", mediaHandler.HandleMediaNew)
	srv.POST("/media", mediaHandler.HandleMediaCreate)
	srv.GET("/media", mediaHandler.HandleMediaShowAll)

	srv.Static("/files", "data/files")

	srv.Start(fmt.Sprintf(":%s", conf.SrvPort))
}
