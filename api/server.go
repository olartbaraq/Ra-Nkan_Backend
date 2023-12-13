package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type Server struct {
	queries *db.Queries
	router  *gin.Engine
	config  *utils.Config
}

var tokenManager *utils.JWTToken

func NewServer(envPath string) *Server {

	config, err := utils.LoadDBConfig(envPath)
	if err != nil {
		panic(fmt.Sprintf("Could not load env config: %v", err))
	}

	conn, err := sql.Open(config.DBdriver, config.DBsourceLive)
	if err != nil {
		panic(fmt.Sprintf("There was an error connecting to database: %v", err))
	}

	tokenManager = utils.NewJWTToken(config)

	q := db.New(conn)

	g := gin.Default()

	g.Use(cors.Default())
	return &Server{
		queries: q,
		router:  g,
		config:  config,
	}

}

func (s *Server) Start(port int) {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		v.RegisterValidation("passwordStrength", ValidatePassword)
		v.RegisterValidation("isImageURL", ImageURLValidation)
	}

	s.router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"Home": "Welcome to Ra'Nkan Homepage...",
		})
	})

	User{}.router(s)
	Auth{}.router(s)
	Category{}.router(s)
	SubCategory{}.router(s)
	Shop{}.router(s)
	Product{}.router(s)

	s.router.Run(fmt.Sprintf(":%d", port))
}
