package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)
	fmt.Println(">> before:", accountOne.Balance, accountTwo.Balance)

	// run a concurrent transfer transaction
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
				result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: accountOne.ID,
				ToAccountID:   accountTwo.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	//check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, accountOne.ID, transfer.FromAccountID)
		require.Equal(t, accountTwo.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, accountOne.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, accountTwo.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//Check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, accountOne.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, accountTwo.ID, toAccount.ID)

		//check balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diffOne := accountOne.Balance - fromAccount.Balance
		diffTwo := toAccount.Balance - accountTwo.Balance
		require.Equal(t, diffOne, diffTwo)
		require.True(t, diffOne > 0)
		require.True(t, diffOne%amount == 0)

		
		k := int(diffOne / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}
	// check the final updated balance
	updatedAccountOne, err := testQueries.GetAccount(context.Background(), accountOne.ID)
	require.NoError(t, err)

	updatedAccountTwo, err := testQueries.GetAccount(context.Background(), accountTwo.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccountOne.Balance, updatedAccountTwo.Balance) 

	require.Equal(t, accountOne.Balance-int64(n)*amount, updatedAccountOne.Balance)
	require.Equal(t, accountTwo.Balance+int64(n)*amount, updatedAccountTwo.Balance)

}




func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)
	fmt.Println(">> before:", accountOne.Balance, accountTwo.Balance)

	// run a concurrent transfer transaction
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := accountOne.ID
		toAccountID := accountTwo.ID

		if i % 2 == 1 {
			fromAccountID = accountTwo.ID
			toAccountID = accountOne.ID
		}

		go func() {
				
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	//check results

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccountOne, err := testQueries.GetAccount(context.Background(), accountOne.ID)
	require.NoError(t, err)

	updatedAccountTwo, err := testQueries.GetAccount(context.Background(), accountTwo.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccountOne.Balance, updatedAccountTwo.Balance) 

	require.Equal(t, accountOne.Balance, updatedAccountOne.Balance)
	require.Equal(t, accountTwo.Balance, updatedAccountTwo.Balance)

}
