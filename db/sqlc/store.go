package db

import (
	"context"
	"database/sql"
	"fmt"
)

// store provides all funcs to execute db queries and txs
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return nil
	}

	// create new query obj for func to perform actions
	q := New(tx)
	// execute the func
	err = fn(q)
	// if err then attempt rollback
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccount int64 `json:"from_account"`
	ToAccount   int64 `json:"to_account"`
	Amount      int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		// Create transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccount: arg.FromAccount,
			ToAccount:   arg.ToAccount,
			Amount:      arg.Amount,
		})

		if err != nil {
			return err
		}

		// Create sender entry record
		fmt.Println(txName, "create sender entry")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccount,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		// Create receiver entry record
		fmt.Println(txName, "create receiver entry")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccount,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// TODO: prevent deadlock and update accounts' balance
		// get account
		fmt.Println(txName, "get sender account")
		account1, err := q.GetAccountForUpdate(context.Background(), arg.FromAccount)
		result.FromAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
			ID:      arg.FromAccount,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "get receiver account")
		account2, err := q.GetAccountForUpdate(context.Background(), arg.ToAccount)
		result.ToAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
			ID:      arg.ToAccount,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
