package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomTransfer(t *testing.T, accountOne Account, accountTwo Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: accountOne.ID,
		ToAccountID:   accountTwo.ID,
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

func TestCreateRandomTransfer(t *testing.T) {
	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)
	createRandomTransfer(t, accountOne, accountTwo)
}

func TestGetTransfer(t *testing.T) {
	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)
	transferOne := createRandomTransfer(t, accountOne, accountTwo)
	transferTwo, err := testQueries.GetTransfer(context.Background(), transferOne.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transferTwo)

	require.Equal(t, transferOne.ID, transferTwo.ID)
	require.Equal(t, transferOne.FromAccountID, transferTwo.FromAccountID)
	require.Equal(t, transferOne.ToAccountID, transferTwo.ToAccountID)
	require.Equal(t, transferOne.Amount, transferTwo.Amount)
	require.WithinDuration(t, transferOne.CreatedAt, transferTwo.CreatedAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, accountOne, accountTwo)
		createRandomTransfer(t, accountTwo, accountOne)
	}
	arg := ListTransfersParams{
		FromAccountID: accountOne.ID,
		ToAccountID: accountTwo.ID,
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t,transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
		require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	}
}
