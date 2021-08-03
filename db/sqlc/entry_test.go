package db

import (
	"context"
	"github.com/jwambugu/go-simple-bank-class/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomEntry(t *testing.T, a Account) Entry {
	arg := CreateEntryParams{
		AccountID: int64(a.ID),
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestQueries_CreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestQueries_GetEntry(t *testing.T) {
	account := createRandomAccount(t)
	newEntry := createRandomEntry(t, account)

	fetchedEntry, err := testQueries.GetEntry(context.Background(), newEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, fetchedEntry)

	require.Equal(t, newEntry.ID, fetchedEntry.ID)
	require.Equal(t, newEntry.AccountID, fetchedEntry.AccountID)
	require.Equal(t, newEntry.Amount, fetchedEntry.Amount)
	require.WithinDuration(t, newEntry.CreatedAt, fetchedEntry.CreatedAt, time.Second)
}

func TestQueries_ListEntries(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		AccountID: int64(account.ID),
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
