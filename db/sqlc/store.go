package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct{
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store{
	return &Store{
		db: db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction 
func(store *Store) execTx(ctx context.Context, fn func(*Queries) error) error{
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil{
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil{
		if rbErr := tx.Rollback(); rbErr != nil{
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams holds the input parameters of the transfer transaction.
type TransferTxParams struct{
	FromAccountID   int64 `json:"from_account_id"`
	ToAccountID    int64 `json:"to_account_id"`
	Amount int64  `json:"amount"`
}

// TransferTxResult contains the result of the transaction.
type TransferTxResult struct{
	Transfer Transfer  `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount Account  `json:"to_account"`
	FromEntry Entry  `json:"from_entry"`
	ToEntry  Entry `json:"to_entry"`
}

// Transfer performs a money transfer from one account to another.
// it creates a transfer record, add account entries and update the account balance within a single database transaction.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error){
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// create transfer record
		result.Transfer, err = q.CreateTransfer(ctx,CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}
		// from entry
		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil{
			return err
		}

		// to entry
		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil{
			return err
		}

		

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		return nil
	})
	return result, err
}
