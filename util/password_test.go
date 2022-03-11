package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	require.NoError(t, CheckPassword(hashedPassword, password))

	wrongPassword := RandomString(7)

	require.EqualError(t, CheckPassword(hashedPassword, wrongPassword), bcrypt.ErrMismatchedHashAndPassword.Error())
}
