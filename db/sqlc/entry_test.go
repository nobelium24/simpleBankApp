package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount: util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, arg.AccountID, entry.AccountID)

	require.NotZero(t, entry.AccountID)
	require.NotZero(t, entry.Amount)

	return entry
}

func TestCreateRandomAccount(t *testing.T){
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t* testing.T){
	account := createRandomAccount(t)
	entryOne := createRandomEntry(t, account)
	entryTwo, err := testQueries.GetEntry(context.Background(), entryOne.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryTwo)

	require.Equal(t, entryOne.ID, entryTwo.ID)
	require.Equal(t, entryOne.AccountID, entryTwo.AccountID)
	require.Equal(t, entryOne.Amount, entryTwo.Amount)
	require.WithinDuration(t, entryOne.CreatedAt, entryTwo.CreatedAt, time.Second)
}

func TestListEntry(t *testing.T){
	account := createRandomAccount(t)
	for i := 0; i < 10; i++{
		createRandomEntry(t, account)
	}
	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit: 5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}