package api

import (
	"context"
	"net/http"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type User struct {
	server *Server
}

type UserParams struct {
	Lastname  string `json:"lastname" binding:"required"`
	Firstname string `json:"firstname" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required,len=11"`
	Address   string `json:"address" binding:"required"`
	Password  string `json:"password" binding:"required,min=8,passwordStrength"`
	IsAdmin   bool   `json:"is_admin"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Lastname  string    `json:"lastname"`
	Firstname string    `json:"firstname"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type DeleteUserParam struct {
	ID int64 `json:"id"`
}

func (u User) router(server *Server) {
	u.server = server
	serverGroup := server.router.Group("/users")
	serverGroup.GET("/allUsers", u.ListUsers)
	serverGroup.POST("/register", u.createUser)
	serverGroup.DELETE("/deactivate", u.DeleteUser)
}

// ValidatePassword checks if the password meets the specified criteria.
func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check if the password is at least 8 characters long
	if utf8.RuneCountInString(password) < 8 {
		return false
	}

	// Check if the password contains at least one digit and one symbol
	hasDigit := false
	hasSymbol := false
	hasUpper := false
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasDigit = true
		} else if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSymbol = true
		} else if unicode.IsUpper(char) {
			hasUpper = true
		}
	}

	return hasDigit && hasSymbol && hasUpper
}

// Register the custom validation function
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("passwordStrength", ValidatePassword)
	}
}

func (u *User) createUser(ctx *gin.Context) {
	user := UserParams{}

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
		Email:          user.Email,
		Phone:          user.Phone,
		Address:        user.Address,
		IsAdmin:        user.IsAdmin,
		HashedPassword: hashedPassword,
	}

	userToSave, err := u.server.queries.CreateUser(context.Background(), arg)

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

func (u *User) ListUsers(ctx *gin.Context) {

	arg := db.ListAllUsersParams{
		Limit:  10,
		Offset: 0,
	}

	users, err := u.server.queries.ListAllUsers(context.Background(), arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	allUsers := []UserResponse{}

	for _, v := range users {

		userResponse := UserResponse{
			ID:        v.ID,
			Lastname:  v.Lastname,
			Firstname: v.Firstname,
			Email:     v.Email,
			Phone:     v.Phone,
			Address:   v.Address,
			IsAdmin:   v.IsAdmin,
			CreatedAt: v.CreatedAt,
		}
		n := userResponse
		allUsers = append(allUsers, n)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "users fetched sucessfully",
		"data":    allUsers,
	})
}

func (u *User) DeleteUser(ctx *gin.Context) {

	id := DeleteUserParam{}

	if err := ctx.ShouldBindJSON(&id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	err := u.server.queries.DeleteUser(context.Background(), id.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "user deactivated sucessfully",
	})
}
