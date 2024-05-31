package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	secretKey := "Th!sIsYour@ppKey&~"
	plaintext := "Hello, World!"

	c := NewCryptWithKey(secretKey)
	encrypted, err := c.Encrypt(plaintext)
	assert.NoError(t, err, "Error encrypting")

	decrypted, err := c.Decrypt(encrypted)
	assert.NoError(t, err, "Error decrypting")

	assert.Equal(t, plaintext, decrypted, "Decrypted text does not match the original text")
}

func TestHashStr(t *testing.T) {
	str := "Hello, World!"
	hashedStr := HashStr(str)
	assert.NotEqual(t, str, hashedStr, "Hashed string should not be the same as the original string")
	assert.Equal(t, len(hashedStr), 64, "SHA-256 hash should be 64 characters long")

	hashedStr2 := HashStr(str)
	assert.Equal(t, hashedStr, hashedStr2, "Hashed string should be the same for the same input")
}

func TestHashPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err, "Error hashing password")
	assert.NotEqual(t, password, hashedPassword, "Hashed password should not be the same as the original password")
}

func TestCheckPassword(t *testing.T) {
	password := "password123"
	hashedPassword, _ := HashPassword(password)
	isValid := CheckPassword(hashedPassword, password)
	assert.True(t, isValid, "Password should be valid")

	isValid = CheckPassword(hashedPassword, "wrongpassword")
	assert.False(t, isValid, "Password should not be valid")
}
