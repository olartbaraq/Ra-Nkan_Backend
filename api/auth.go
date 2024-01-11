package api

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type Auth struct {
	server *Server
}

type CreateUserParams struct {
	Lastname  string `json:"lastname" binding:"required"`
	Firstname string `json:"firstname" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required,len=11"`
	Address   string `json:"address" binding:"required"`
	Password  string `json:"password" binding:"required,passwordStrength"`
	IsAdmin   bool   `json:"is_admin"`
}

type LoginUserParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (a Auth) router(server *Server) {

	a.server = server

	serverGroup := server.router.Group("/auth")
	serverGroup.POST("/register", a.register)
	serverGroup.POST("/login", a.login)
	serverGroup.GET("/refresh", a.RefreshAccessToken)
	serverGroup.GET("/logout", a.LogoutUser)
}

func (a *Auth) register(ctx *gin.Context) {

	passwordStrengthResp := []string{
		"Password must be minimum of 8 characters",
		"Password must contain at least a number",
		"Password must contain at least a symbol",
		"Password must contain an upper case letter",
		"Password must contain a lower case letter",
	}

	user := CreateUserParams{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		//fmt.Println(err.Error())
		stringErr := string(err.Error())
		//fmt.Println(stringErr)
		if strings.Contains(stringErr, "passwordStrength") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "password Strength not met",
				"Error":   passwordStrengthResp,
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	hashedPassword, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	arg := db.CreateUserParams{
		Lastname:       user.Lastname,
		Firstname:      user.Firstname,
		Email:          strings.ToLower(user.Email),
		Phone:          user.Phone,
		Address:        user.Address,
		IsAdmin:        user.IsAdmin,
		HashedPassword: hashedPassword,
	}

	userToSave, err := a.server.queries.CreateUser(context.Background(), arg)

	if err != nil {
		handleCreateUserError(ctx, err)
		return
	}

	userResponse := UserResponse{
		ID:         userToSave.ID,
		Lastname:   userToSave.Lastname,
		Firstname:  userToSave.Firstname,
		Email:      userToSave.Email,
		Phone:      userToSave.Phone,
		Address:    userToSave.Address,
		IsAdmin:    userToSave.IsAdmin,
		IsLoggedIn: "false",
		CreatedAt:  userToSave.CreatedAt,
		UpdatedAt:  userToSave.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"status":     "success",
		"message":    "user created successfully",
		"data":       userResponse,
	})
}

func handleCreateUserError(ctx *gin.Context, err error) {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			// to check for unique constraint
			handleUniqueConstraintError(ctx, pqErr)
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
		}
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
	}
}

func handleUniqueConstraintError(ctx *gin.Context, pqErr *pq.Error) {
	stringErr := string(pqErr.Detail)
	if strings.Contains(stringErr, "phone") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "User with phone number already exists",
		})
	} else if strings.Contains(stringErr, "email") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "User with email address already exists",
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "Duplicate key violation",
		})
	}
}

func (a Auth) login(ctx *gin.Context) {
	userToLogin := LoginUserParams{}

	if err := ctx.ShouldBindJSON(&userToLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	dbUser, err := a.server.queries.GetUserByEmail(context.Background(), strings.ToLower(userToLogin.Email))

	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"Error":      err.Error(),
			"message":    "The requested user with the specified email does not exist.",
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

	access_token, err := tokenManager.CreateToken(dbUser.ID, dbUser.IsAdmin, a.server.config2.AccessTokenExpiresIn)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	refresh_token, err := tokenManager.CreateToken(dbUser.ID, dbUser.IsAdmin, a.server.config2.RefreshTokenExpiresIn)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.SetCookie("access_token", access_token, a.server.config2.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, a.server.config2.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", a.server.config2.AccessTokenMaxAge*60, "/", "localhost", false, false)

	userResponse := UserResponse{
		ID:         dbUser.ID,
		Lastname:   dbUser.Lastname,
		Firstname:  dbUser.Firstname,
		Email:      dbUser.Email,
		Phone:      dbUser.Phone,
		Address:    dbUser.Address,
		IsAdmin:    dbUser.IsAdmin,
		IsLoggedIn: "true",
		CreatedAt:  dbUser.CreatedAt,
		UpdatedAt:  dbUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"status":     "success",
		"message":    "login successful",
		//"access_token": access_token,
		//"refresh_token": refresh_token,
		"data": userResponse,
	})
}

func (a *Auth) RefreshAccessToken(ctx *gin.Context) {

	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	userId, _, err := tokenManager.VerifyToken(cookie)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := a.server.queries.GetUserById(context.Background(), userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	access_token, err := tokenManager.CreateToken(user.ID, user.IsAdmin, a.server.config2.AccessTokenExpiresIn)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, a.server.config2.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", a.server.config2.AccessTokenMaxAge*60, "/", "localhost", false, false)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func (a *Auth) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
