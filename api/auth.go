package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type Auth struct {
	server *Server
}

type LoginUserParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (a Auth) router(server *Server) {

	a.server = server

	serverGroup := server.router.Group("/auth")
	serverGroup.POST("/login", a.login)
}

func (a Auth) login(ctx *gin.Context) {
	userToLogin := LoginUserParams{}

	if err := ctx.ShouldBindJSON(&userToLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	dbUser, err := a.server.queries.GetUserByEmail(context.Background(), userToLogin.Email)

	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Error":   err.Error(),
			"message": "The requested user with the specified email does not exist.",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	err = utils.VerifyPassword(userToLogin.Password, dbUser.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"Error":   err.Error(),
			"message": "Invalid password. Please check your credentials and try again.",
		})
		return
	}

	token, err := utils.CreateToken(dbUser.ID, dbUser.IsAdmin, a.server.config.SigningKey)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	userResponse := UserResponse{
		ID:        dbUser.ID,
		Lastname:  dbUser.Lastname,
		Firstname: dbUser.Firstname,
		Email:     dbUser.Email,
		Phone:     dbUser.Phone,
		Address:   dbUser.Address,
		IsAdmin:   dbUser.IsAdmin,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "login successful",
		"token":   token,
		"data":    userResponse,
	})
}
