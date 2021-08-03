package db

import (
	"context"
	"github.com/jwambugu/go-simple-bank-class/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T, a, b Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: int64(a.ID),
		ToAccountID:   int64(b.ID),
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestQueries_CreateTransfer(t *testing.T) {
	a := createRandomAccount(t)
	b := createRandomAccount(t)

	createRandomTransfer(t, a, b)
}

func TestQueries_GetTransfer(t *testing.T) {
	a := createRandomAccount(t)
	b := createRandomAccount(t)

	newTransfer := createRandomTransfer(t, a, b)

	foundTransfer, err := testQueries.GetTransfer(context.Background(), newTransfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, foundTransfer)

	require.Equal(t, newTransfer.ID, foundTransfer.ID)
	require.Equal(t, newTransfer.FromAccountID, foundTransfer.FromAccountID)
	require.Equal(t, newTransfer.ToAccountID, foundTransfer.ToAccountID)
	require.Equal(t, newTransfer.Amount, foundTransfer.Amount)
	require.WithinDuration(t, newTransfer.CreatedAt, foundTransfer.CreatedAt, time.Second)
}

func TestQueries_ListTransfers(t *testing.T) {
	a := createRandomAccount(t)
	b := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, a, b)
		createRandomTransfer(t, a, b)
	}

	arg := ListTransfersParams{
		FromAccountID: int64(a.ID),
		ToAccountID:   int64(a.ID),
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == int64(a.ID) && transfer.ToAccountID == int64(b.ID))
	}
}
