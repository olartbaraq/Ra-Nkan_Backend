package all_test

import (
	"fmt"
	"testing"

	"github.com/olartbaraq/spectrumshelf/utils"
	"github.com/stretchr/testify/assert"
)

type Params struct {
	userID  int64
	isAdmin bool
}

func returnJWT() *utils.JWTToken {
	config, err := utils.LoadDBConfig("..")
	if err != nil {
		panic(fmt.Sprintf("Could not load env config: %v", err))
	}

	jwtToken := utils.NewJWTToken(config)

	return jwtToken
}

func TestCreateToken(t *testing.T) {

	userToken := Params{
		userID:  1,
		isAdmin: false,
	}

	token, err := returnJWT().CreateToken(userToken.userID, userToken.isAdmin)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestVerifyToken(t *testing.T) {

	userToken := Params{
		userID:  10,
		isAdmin: true,
	}

	token, err := returnJWT().CreateToken(userToken.userID, userToken.isAdmin)

	if err != nil {
		t.Fatalf("Error generating token: %v", err)
		return
	}

	claimToken, role, err := (returnJWT()).VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, role, "admin")
	assert.NotEmpty(t, claimToken)
	assert.NotZero(t, claimToken)
	assert.Equal(t, claimToken, userToken.userID)
}
