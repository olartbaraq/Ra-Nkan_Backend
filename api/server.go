package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type Server struct {
	queries *db.Queries
	router  *gin.Engine
	config  *utils.Config
}

func NewServer(envPath string) *Server {

	config, err := utils.LoadConfig(envPath)
	if err != nil {
		panic(fmt.Sprintf("Could not load env config: %v", err))
	}

	conn, err := sql.Open(config.DBdriver, config.DBsourceLive)
	if err != nil {
		panic(fmt.Sprintf("There was an error connecting to database: %v", err))
	}

	q := db.New(conn)

	g := gin.Default()

	return &Server{
		queries: q,
		router:  g,
		config:  config,
	}

}

func (s *Server) Start(port int) {
	s.router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"Home": "Welcome to Ra'Nkan Homepage...",
		})
	})

	User{}.router(s)
	Auth{}.router(s)

	s.router.Run(fmt.Sprintf(":%d", port))
}
