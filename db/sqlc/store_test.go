package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)

	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)

	// Cast the ids
	//accountOne.ID = int32(accountOne.ID)

	// run n concurrent transfer transaction
	n := 5
	amount := int64(10)

	errsChan := make(chan error)
	resultChan := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()

			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: int64(accountOne.ID),
				ToAccountID:   int64(accountTwo.ID),
				Amount:        amount,
			})

			errsChan <- err
			resultChan <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errsChan
		require.NoError(t, err)

		result := <-resultChan
		require.NotEmpty(t, result)

		// Check the result transfer
		transfer := result.Transfer

		require.NotEmpty(t, transfer)
		require.Equal(t, int64(accountOne.ID), transfer.FromAccountID)
		require.Equal(t, int64(accountTwo.ID), transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check the sender entry
		fromEntry := result.FromEntry

		require.NotEmpty(t, fromEntry)
		require.Equal(t, int64(accountOne.ID), fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// Check the receiver entry
		toEntry := result.ToEntry

		require.NotEmpty(t, toEntry)
		require.Equal(t, int64(accountTwo.ID), toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO: check account balance
	}
}
