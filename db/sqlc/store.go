package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)

	if err := fn(q); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rollbackErr)
		}

		return err
	}

	return tx.Commit()
}

func addMoney(ctx context.Context, q *Queries, fromAccountID, fromAmount, toAccountID, toAmount int64) (
	from, to Account, err error) {

	from, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     int32(fromAccountID),
		Amount: fromAmount,
	})

	if err != nil {
		return
	}

	to, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     int32(toAccountID),
		Amount: toAmount,
	})

	return
}

// TransferTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create a new transfer between the accounts
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		// Debit money from the sender
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		// Credit money to the receiver
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			// Update the sender's account balance first
			result.FromAccount, result.ToAccount, err = addMoney(context.Background(), q, arg.FromAccountID, -arg.Amount,
				arg.ToAccountID, arg.Amount)
		} else {
			// Update the receiver's account balance first
			result.ToAccount, result.FromAccount, err = addMoney(context.Background(), q, arg.ToAccountID, arg.Amount,
				arg.FromAccountID, -arg.Amount)
		}

		return nil
	})

	return result, err
}
