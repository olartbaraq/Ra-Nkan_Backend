package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
)

type Oauth struct {
	server *Server
}

func (o Oauth) router(server *Server) {
	o.server = server

	serverGroup := server.router.Group("/oauth")
	serverGroup.POST("/google/create_user", o.createUser)
}

type OauthParams struct {
	ID_token string `json:"id_token" binding:"required"`
}

type UserInfo struct {
	FamilyName     string `json:"family_name" binding:"required"`
	GivenName      string `json:"given_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	IsAdmin        bool   `json:"is_admin"`
	HashedPassword string `json:"hashed_password"`
	// Picture       string `json:"picture"`
	// EmailVerified bool   `json:"email_verified"`
	// Locale        string `json:"locale"`
}

func (o *Oauth) createUser(ctx *gin.Context) {

	user_token := OauthParams{}

	if err := ctx.ShouldBindJSON(&user_token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	userInfoChan := make(chan *UserInfo, 1)

	go func(user_token string, u chan<- *UserInfo) {
		url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", user_token)

		resp, err := http.Get(url)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error":   err.Error(),
				"message": "error making request to Google API",
			})
			return

		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error":   err.Error(),
				"message": "google API returned non-OK status",
			})
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error":   err.Error(),
				"message": "error reading response body",
			})
			return
		}

		var userInfo UserInfo
		if err := json.Unmarshal(body, &userInfo); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error":   err.Error(),
				"message": "error unmarshalling JSON response",
			})
			return
		}

		u <- &userInfo
		//fmt.Println(userInfo)
	}(user_token.ID_token, userInfoChan)

	select {
	case userInfo := <-userInfoChan:

		oAuthDbUser, err := o.server.queries.GetUserByEmail(context.Background(), strings.ToLower(userInfo.Email))

		if err == sql.ErrNoRows {

			arg := db.CreateUserParams{
				Lastname:       userInfo.FamilyName,
				Firstname:      userInfo.GivenName,
				Email:          strings.ToLower(userInfo.Email),
				Phone:          "",
				Address:        "",
				HashedPassword: "",
			}

			oAuthUserToSave, err := o.server.queries.CreateUser(context.Background(), arg)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"Error": err.Error(),
				})
				return
			}

			token, err := tokenManager.CreateToken(oAuthUserToSave.ID, oAuthUserToSave.IsAdmin)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"Error": err.Error(),
				})
				return
			}

			userResponse := UserResponse{
				ID:        oAuthUserToSave.ID,
				Lastname:  oAuthUserToSave.Lastname,
				Firstname: oAuthUserToSave.Firstname,
				Email:     oAuthUserToSave.Email,
				Phone:     oAuthUserToSave.Phone,
				Address:   oAuthUserToSave.Address,
				IsAdmin:   oAuthUserToSave.IsAdmin,
				CreatedAt: oAuthUserToSave.CreatedAt,
				UpdatedAt: oAuthUserToSave.UpdatedAt,
			}

			ctx.JSON(http.StatusOK, gin.H{
				"statusCode": http.StatusOK,
				"status":     "success",
				"message":    "login successful",
				"token":      token,
				"data":       userResponse,
			})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			return
		}

		token, err := tokenManager.CreateToken(oAuthDbUser.ID, oAuthDbUser.IsAdmin)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			return
		}

		userResponse := UserResponse{
			ID:        oAuthDbUser.ID,
			Lastname:  oAuthDbUser.Lastname,
			Firstname: oAuthDbUser.Firstname,
			Email:     oAuthDbUser.Email,
			Phone:     oAuthDbUser.Phone,
			Address:   oAuthDbUser.Address,
			IsAdmin:   oAuthDbUser.IsAdmin,
			CreatedAt: oAuthDbUser.CreatedAt,
			UpdatedAt: oAuthDbUser.UpdatedAt,
		}

		ctx.JSON(http.StatusOK, gin.H{
			"statusCode": http.StatusOK,
			"status":     "success",
			"message":    "login successful",
			"token":      token,
			"data":       userResponse,
		})

	case <-time.After(10 * time.Second):
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "timeout fetching user information",
		})
		return
	}

	close(userInfoChan)

}
