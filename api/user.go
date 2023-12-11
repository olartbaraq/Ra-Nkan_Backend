package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type User struct {
	server *Server
}

type CreateUserParams struct {
	Lastname  string `json:"lastname" binding:"required"`
	Firstname string `json:"firstname" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required,len=11"`
	Address   string `json:"address" binding:"required"`
	Password  string `json:"password" binding:"required,min=8" validate:"passwordStrength"`
	IsAdmin   bool   `json:"is_admin"`
}

type UpdateUserParams struct {
	ID        int64     `json:"id" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Phone     string    `json:"phone" binding:"required,len=11"`
	Address   string    `json:"address" binding:"required"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserPasswordParams struct {
	ID        int64     `json:"id" binding:"required"`
	Password  string    `json:"password" binding:"required,min=8" validate:"passwordStrength"`
	UpdatedAt time.Time `json:"updated_at"`
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
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteUserParam struct {
	ID int64 `json:"id"`
}

func (u User) router(server *Server) {
	u.server = server
	serverGroup := server.router.Group("/users", AuthenticatedMiddleware())
	serverGroup.GET("/allUsers", u.listUsers)
	serverGroup.PUT("/update", u.updateUser)
	serverGroup.PUT("/update/password", u.updatePassword)
	serverGroup.DELETE("/deactivate", u.deleteUser)
	serverGroup.GET("/profile", u.userProfile)
}

func extractTokenFromRequest(ctx *gin.Context) (string, error) {
	// Extract the token from the Authorization header
	authorizationHeader := ctx.GetHeader("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("unauthorized request")
	}

	// Expecting the header to be in the format "Bearer <token>"
	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 && strings.ToLower(headerParts[0]) != "bearer" {
		return "", errors.New("invalid token format")
	}

	return headerParts[1], nil
}

func returnIdRole(tokenString string) (int64, string, error) {

	if tokenString == "" {
		return 0, "", errors.New("unauthorized: Missing or invalid token")
	}

	userId, role, err := tokenManager.VerifyToken(tokenString)

	if err != nil {
		return 0, "", errors.New("failed to verify token")
	}

	return userId, role, nil
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

func (u *User) listUsers(ctx *gin.Context) {

	tokenString, err := extractTokenFromRequest(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Missing or invalid token",
		})
		return
	}

	_, role, err := returnIdRole(tokenString)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":  err.Error(),
			"status": "failed to verify token",
		})
		ctx.Abort()
		return
	}

	if role != utils.AdminRole {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		ctx.Abort()
		return
	}

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
		"message": "all users fetched sucessfully",
		"data":    allUsers,
	})
}

func (u *User) deleteUser(ctx *gin.Context) {

	tokenString, err := extractTokenFromRequest(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Missing or invalid token",
		})
		ctx.Abort()
		return
	}

	userId, _, err := returnIdRole(tokenString)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":  err.Error(),
			"status": "failed to verify token",
		})
		ctx.Abort()
		return
	}

	id := DeleteUserParam{}

	if userId != id.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: invalid token",
		})
		ctx.Abort()
		return
	}

	if err := ctx.ShouldBindJSON(&id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	err = u.server.queries.DeleteUser(context.Background(), id.ID)

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

func (u *User) updateUser(ctx *gin.Context) {

	tokenString, err := extractTokenFromRequest(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Missing or invalid token",
		})
		return
	}

	userId, _, err := returnIdRole(tokenString)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":  err.Error(),
			"status": "failed to verify token",
		})
		ctx.Abort()
		return
	}

	user := UpdateUserParams{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	if userId != user.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: invalid token",
		})
		ctx.Abort()
		return
	}

	arg := db.UpdateUserParams{
		ID:        user.ID,
		Email:     strings.ToLower(user.Email),
		Phone:     user.Phone,
		Address:   user.Address,
		UpdatedAt: time.Now(),
	}

	userToUpdate, err := u.server.queries.UpdateUser(context.Background(), arg)

	if err != nil {
		handleCreateUserError(ctx, err)
		return
	}

	userResponse := UserResponse{
		ID:        userToUpdate.ID,
		Lastname:  userToUpdate.Lastname,
		Firstname: userToUpdate.Firstname,
		Email:     userToUpdate.Email,
		Phone:     userToUpdate.Phone,
		Address:   userToUpdate.Address,
		IsAdmin:   userToUpdate.IsAdmin,
		CreatedAt: userToUpdate.CreatedAt,
		UpdatedAt: userToUpdate.UpdatedAt,
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "user updated successfully",
		"data":    userResponse,
	})
}

func (u *User) userProfile(ctx *gin.Context) {
	value, exist := ctx.Get("id")

	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  exist,
			"message": "Unauthorized",
		})
		return
	}

	userId, ok := value.(int64)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  exist,
			"message": "Issue Encountered, try again later",
		})
		return
	}

	user, err := u.server.queries.GetUserById(context.Background(), userId)

	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Error":   err.Error(),
			"message": "Unauthorized",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Issue Encountered, try again later",
		})
		return
	}

	userResponse := UserResponse{
		ID:        user.ID,
		Lastname:  user.Lastname,
		Firstname: user.Firstname,
		Email:     user.Email,
		Phone:     user.Phone,
		Address:   user.Address,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user fetched successfully",
		"data":    userResponse,
	})
}

func (u *User) updatePassword(ctx *gin.Context) {

	tokenString, err := extractTokenFromRequest(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Missing or invalid token",
		})
		return
	}

	userId, _, err := returnIdRole(tokenString)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":  err.Error(),
			"status": "failed to verify token",
		})
		ctx.Abort()
		return
	}

	user := UpdateUserPasswordParams{}

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
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	if userId != user.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: invalid token",
		})
		ctx.Abort()
		return
	}

	hashedPassword, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	arg := db.UpdateUserPasswordParams{
		ID:             user.ID,
		HashedPassword: hashedPassword,
		UpdatedAt:      time.Now(),
	}

	userToUpdatePassword, err := u.server.queries.UpdateUserPassword(context.Background(), arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	userResponse := UserResponse{
		ID:        userToUpdatePassword.ID,
		Lastname:  userToUpdatePassword.Lastname,
		Firstname: userToUpdatePassword.Firstname,
		Email:     userToUpdatePassword.Email,
		Phone:     userToUpdatePassword.Phone,
		Address:   userToUpdatePassword.Address,
		IsAdmin:   userToUpdatePassword.IsAdmin,
		CreatedAt: userToUpdatePassword.CreatedAt,
		UpdatedAt: userToUpdatePassword.UpdatedAt,
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "password updated successfully",
		"data":    userResponse,
	})
}
