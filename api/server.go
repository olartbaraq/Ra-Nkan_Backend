package api

import (
	"database/sql"
	_ "log"
	_ "net/http"

	"github.com/gin-gonic/gin"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type Server struct {
	queries *db.Queries
	router  *gin.Engine
}

func NewServer() *Server {

	config, err := utils.LoadConfig("../..")
	if err != nil {
		panic("Could not load env config:")
	}

	conn, err := sql.Open(config.DBdriver, config.DBsource)
	if err != nil {
		panic("There was an error connecting to database")
	}

	q := db.New(conn)

	g := gin.Default()

	return &Server{
		queries: q,
		router:  g,
	}
}
