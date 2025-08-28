package services

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hashed, err := hashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	assert.NoError(t, err)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	result := checkPasswordHash(password, string(hashed))
	assert.True(t, result)

	result = checkPasswordHash("wrongpassword", string(hashed))
	assert.False(t, result)
}
