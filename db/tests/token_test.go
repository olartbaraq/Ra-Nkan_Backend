package db_test

import (
	"testing"

	"github.com/olartbaraq/spectrumshelf/utils"
	"github.com/stretchr/testify/assert"
)

type Params struct {
	userID     int64
	isAdmin    bool
	signingKey string
}

func TestCreateToken(t *testing.T) {

	userToken := Params{
		userID:     1,
		isAdmin:    false,
		signingKey: "1234567890rt569dhgk90565wdjgjlyu",
	}

	token, err := utils.CreateToken(userToken.userID, userToken.isAdmin, userToken.signingKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestVerifyToken(t *testing.T) {

	userToken := Params{
		userID:     10,
		isAdmin:    true,
		signingKey: "1234567890rt569dhgk90565wdjg6748",
	}

	token, err := utils.CreateToken(userToken.userID, userToken.isAdmin, userToken.signingKey)

	if err != nil {
		t.Fatalf("Error generating token: %v", err)
		return
	}

	claimToken, role, err := utils.VerifyToken(token, userToken.signingKey)
	assert.NoError(t, err)
	assert.Equal(t, role, "admin")
	assert.NotEmpty(t, claimToken)
	assert.NotZero(t, claimToken)
	assert.Equal(t, claimToken, userToken.userID)
}
