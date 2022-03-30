package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/colincassens/cosmosCheckers/testutil/keeper"
	"github.com/colincassens/cosmosCheckers/testutil/nullify"
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/keeper"
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/types"
)

func createTestNextGame(keeper *keeper.Keeper, ctx sdk.Context) types.NextGame {
	item := types.NextGame{}
	keeper.SetNextGame(ctx, item)
	return item
}

func TestNextGameGet(t *testing.T) {
	keeper, ctx := keepertest.CosmoscheckersKeeper(t)
	item := createTestNextGame(keeper, ctx)
	rst, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestNextGameRemove(t *testing.T) {
	keeper, ctx := keepertest.CosmoscheckersKeeper(t)
	createTestNextGame(keeper, ctx)
	keeper.RemoveNextGame(ctx)
	_, found := keeper.GetNextGame(ctx)
	require.False(t, found)
}
