package db

import (
	"context"
	"github.com/jwambugu/go-simple-bank-class/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		FullName:       util.RandomOwner(),
		HashedPassword: "",
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestQueries_CreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestQueries_GetUser(t *testing.T) {
	expected := createRandomUser(t)

	actual, err := testQueries.GetUser(context.Background(), expected.Username)
	require.NoError(t, err)
	require.NotEmpty(t, actual)

	require.Equal(t, expected.Username, actual.Username)
	require.Equal(t, expected.FullName, actual.FullName)
	require.Equal(t, expected.HashedPassword, actual.HashedPassword)
	require.Equal(t, expected.Email, actual.Email)
	require.WithinDuration(t, expected.PasswordChangedAt, actual.PasswordChangedAt, time.Second)
	require.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Second)
}
