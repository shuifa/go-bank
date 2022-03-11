package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/shuifa/go-bank/db/sqlc"
	"github.com/shuifa/go-bank/token"
	"github.com/shuifa/go-bank/util"
)

// Server serves http request for simple bank project
type Server struct {
	store      db.Store
	tokenMaker token.Maker
	route      *gin.Engine
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("generate tokenMaker err %w", err)
	}

	server := &Server{
		store:      store,
		route:      nil,
		tokenMaker: tokenMaker,
		config:     config,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	route := gin.Default()

	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := validate.RegisterValidation("currency", validateCurrency)
		if err != nil {
			log.Fatal(err)
		}
	}

	route.POST("/users", server.createUser)
	route.POST("/users/login", server.login)

	authRooter := route.Group("/").Use(AuthMiddleware(server.tokenMaker))

	authRooter.POST("/accounts", server.createAccount)
	authRooter.GET("/account/:id", server.getAccount)
	authRooter.GET("/accounts", server.listAccount)

	authRooter.POST("/transfers", server.createTransfer)

	server.route = route
}

func (server *Server) Start(addr string) error {
	return server.route.Run(addr)
}

func ErrResponse(err error) gin.H {
	return gin.H{"error:": err.Error()}
}
