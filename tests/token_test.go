package all_test

import (
	"log"
	"testing"

	"github.com/olartbaraq/spectrumshelf/utils"
	"github.com/stretchr/testify/assert"
)

type Params struct {
	userID  int64
	isAdmin bool
}

var tokenManager *utils.JWTToken

var DbConfig *utils.Config

func TestCreateToken(t *testing.T) {

	DbConfig, err := utils.LoadOtherConfig(".")
	if err != nil {
		log.Fatal("Could not load env config", err)
	}

	userToken := Params{
		userID:  1,
		isAdmin: false,
	}

	tokenManager = utils.NewJWTToken(DbConfig)

	token, err := tokenManager.CreateToken(userToken.userID, userToken.isAdmin, DbConfig.AccessTokenExpiresIn)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestVerifyToken(t *testing.T) {

	DbConfig, err := utils.LoadOtherConfig(".")
	if err != nil {
		log.Fatal("Could not load env config", err)
	}

	userToken := Params{
		userID:  10,
		isAdmin: true,
	}
	tokenManager = utils.NewJWTToken(DbConfig)

	token, err := tokenManager.CreateToken(userToken.userID, userToken.isAdmin, DbConfig.AccessTokenExpiresIn)

	if err != nil {
		t.Fatalf("Error generating token: %v", err)
		return
	}

	claimToken, role, err := tokenManager.VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, role, "admin")
	assert.NotEmpty(t, claimToken)
	assert.NotZero(t, claimToken)
	assert.Equal(t, claimToken, userToken.userID)
}
