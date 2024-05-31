package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenGenerationAndVerification(t *testing.T) {
	// Set the JWT_KEY environment variable for testing
	os.Setenv("JWT_KEY", "test_key")

	// Generate a token for a test user
	userID := "test_user"
	token, err := GenerateToken(userID)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	// Generate a refresh token for the test user
	refreshToken, err := GenerateRefreshToken(userID)
	assert.Nil(t, err)
	assert.NotEmpty(t, refreshToken)

	// Verify the refresh token
	verifiedUserID, err := VerifyRefreshToken(refreshToken)
	assert.Nil(t, err)
	assert.Equal(t, userID, verifiedUserID)
}
