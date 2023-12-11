package api

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"github.com/lib/pq"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
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
	serverGroup.POST("/register", a.register)
	serverGroup.POST("/login", a.login)
}

// Register the custom validation function
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("passwordStrength", ValidatePassword)
	}
}

func (a *Auth) register(ctx *gin.Context) {
	user := CreateUserParams{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		stringErr := string(err.Error())
		if strings.Contains(stringErr, "passwordStrength") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Error": `
						"Password must be minimum of 8 characters",
						"Password must be contain at least a number",
						"Password must be contain at least a symbol",
						"Password must be contain a upper case letter"
						`,
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
		ID:        userToSave.ID,
		Lastname:  userToSave.Lastname,
		Firstname: userToSave.Firstname,
		Email:     userToSave.Email,
		Phone:     userToSave.Phone,
		Address:   userToSave.Address,
		IsAdmin:   userToSave.IsAdmin,
		CreatedAt: userToSave.CreatedAt,
		UpdatedAt: userToSave.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "user created successfully",
		"data":    userResponse,
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

	token, err := tokenManager.CreateToken(dbUser.ID, dbUser.IsAdmin)

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
