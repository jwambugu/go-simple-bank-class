package db

import (
	"context"
	"database/sql"
	"github.com/jwambugu/go-simple-bank-class/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestQueries_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestQueries_GetAccount(t *testing.T) {
	newAccount := createRandomAccount(t)

	fetchedAccount, err := testQueries.GetAccount(context.Background(), newAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, fetchedAccount)

	require.Equal(t, newAccount.ID, fetchedAccount.ID)
	require.Equal(t, newAccount.Owner, fetchedAccount.Owner)
	require.Equal(t, newAccount.Balance, fetchedAccount.Balance)
	require.Equal(t, newAccount.Currency, fetchedAccount.Currency)

	require.WithinDuration(t, newAccount.CreatedAt, fetchedAccount.CreatedAt, time.Second)
}

func TestQueries_UpdateAccount(t *testing.T) {
	newAccount := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      newAccount.ID,
		Balance: util.RandomMoney(),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, newAccount.ID, updatedAccount.ID)
	require.Equal(t, newAccount.Owner, updatedAccount.Owner)
	require.Equal(t, newAccount.Currency, updatedAccount.Currency)
	require.Equal(t, newAccount.Currency, updatedAccount.Currency)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
}

func TestQueries_DeleteAccount(t *testing.T) {
	newAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), newAccount.ID)

	require.NoError(t, err)

	fetchedAccount, err := testQueries.GetAccount(context.Background(), newAccount.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, fetchedAccount)
}

func TestQueries_ListAccounts(t *testing.T) {
	var lastAccount Account

	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
