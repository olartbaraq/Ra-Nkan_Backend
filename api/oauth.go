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
		//fmt.Println(url)

		resp, err := http.Get(url)
		//fmt.Println(resp)
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
				"error": "google API returned non-OK status",
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

		if userInfo.FamilyName != "" && userInfo.GivenName != "" && userInfo.Email != "" {
			u <- &userInfo
		}

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

			access_token, err := tokenManager.CreateToken(oAuthUserToSave.ID, oAuthUserToSave.IsAdmin, o.server.config.AccessTokenExpiresIn)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"Error": err.Error(),
				})
				return
			}

			refresh_token, err := tokenManager.CreateToken(oAuthUserToSave.ID, oAuthUserToSave.IsAdmin, o.server.config.RefreshTokenExpiresIn)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"Error": err.Error(),
				})
				return
			}

			ctx.SetCookie("access_token", access_token, o.server.config.AccessTokenMaxAge*60, "/", "localhost", false, true)
			ctx.SetCookie("refresh_token", refresh_token, ConfigViper.RefreshTokenMaxAge*60, "/", "localhost", false, true)
			ctx.SetCookie("logged_in", "true", o.server.config.AccessTokenMaxAge*60, "/", "localhost", false, false)

			userResponse := UserResponse{
				ID:         oAuthUserToSave.ID,
				Lastname:   oAuthUserToSave.Lastname,
				Firstname:  oAuthUserToSave.Firstname,
				Email:      oAuthUserToSave.Email,
				Phone:      oAuthUserToSave.Phone,
				Address:    oAuthUserToSave.Address,
				IsAdmin:    oAuthUserToSave.IsAdmin,
				IsLoggedIn: "true",
				CreatedAt:  oAuthUserToSave.CreatedAt,
				UpdatedAt:  oAuthUserToSave.UpdatedAt,
			}

			ctx.JSON(http.StatusOK, gin.H{
				"statusCode":    http.StatusOK,
				"status":        "success",
				"message":       "login successful",
				"access_token":  access_token,
				"refresh_token": refresh_token,
				"data":          userResponse,
			})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			return
		}

		access_token, err := tokenManager.CreateToken(oAuthDbUser.ID, oAuthDbUser.IsAdmin, o.server.config.AccessTokenExpiresIn)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			return
		}

		refresh_token, err := tokenManager.CreateToken(oAuthDbUser.ID, oAuthDbUser.IsAdmin, o.server.config.RefreshTokenExpiresIn)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			return
		}

		ctx.SetCookie("access_token", access_token, o.server.config.AccessTokenMaxAge*60, "/", "localhost", false, true)
		ctx.SetCookie("refresh_token", refresh_token, o.server.config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
		ctx.SetCookie("logged_in", "true", o.server.config.AccessTokenMaxAge*60, "/", "localhost", false, false)

		userResponse := UserResponse{
			ID:         oAuthDbUser.ID,
			Lastname:   oAuthDbUser.Lastname,
			Firstname:  oAuthDbUser.Firstname,
			Email:      oAuthDbUser.Email,
			Phone:      oAuthDbUser.Phone,
			Address:    oAuthDbUser.Address,
			IsAdmin:    oAuthDbUser.IsAdmin,
			IsLoggedIn: "true",
			CreatedAt:  oAuthDbUser.CreatedAt,
			UpdatedAt:  oAuthDbUser.UpdatedAt,
		}

		ctx.JSON(http.StatusOK, gin.H{
			"statusCode":    http.StatusOK,
			"status":        "success",
			"message":       "login successful",
			"access_token":  access_token,
			"refresh_token": refresh_token,
			"data":          userResponse,
		})

	case <-time.After(5 * time.Second):
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "timeout fetching user information",
		})
		return
	}

	close(userInfoChan)

}
