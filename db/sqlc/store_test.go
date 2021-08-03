package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)

	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)

	fmt.Println(">> before tx:", accountOne.Balance, accountTwo.Balance)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errsChan := make(chan error)
	resultChan := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)

		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)

			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: int64(accountOne.ID),
				ToAccountID:   int64(accountTwo.ID),
				Amount:        amount,
			})

			errsChan <- err
			resultChan <- result
		}()
	}

	transactionExists := make(map[int]bool)

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

		// Check the sender account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, accountOne.ID, fromAccount.ID)

		// Check the receiver account
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, accountTwo.ID, toAccount.ID)

		// Check account balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		accountOneBalance := accountOne.Balance - fromAccount.Balance
		accountTwoBalance := toAccount.Balance - accountTwo.Balance

		require.Equal(t, accountOneBalance, accountTwoBalance)
		require.True(t, accountOneBalance > 0)
		require.True(t, accountOneBalance%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(accountOneBalance / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, transactionExists, k)
		transactionExists[k] = true
	}

	// Check for the final updated balances
	updatedAccountOne, err := testQueries.GetAccount(context.Background(), accountOne.ID)
	require.NoError(t, err)

	updatedAccountTwo, err := testQueries.GetAccount(context.Background(), accountTwo.ID)
	require.NoError(t, err)

	fmt.Println(">> after tx:", updatedAccountOne.Balance, updatedAccountTwo.Balance)

	// The new account one balance will reduce by n * amount
	require.Equal(t, accountOne.Balance-int64(n)*amount, updatedAccountOne.Balance)

	// The new account two balance will increase by n * amount
	require.Equal(t, accountTwo.Balance+int64(n)*amount, updatedAccountTwo.Balance)
}
