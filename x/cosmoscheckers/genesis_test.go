package cosmoscheckers_test

import (
	"testing"

	keepertest "github.com/colincassens/cosmosCheckers/testutil/keeper"
	"github.com/colincassens/cosmosCheckers/testutil/nullify"
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers"
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		NextGame: &types.NextGame{
			IdValue: 6,
		},
		StoredGameList: []types.StoredGame{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.CosmoscheckersKeeper(t)
	cosmoscheckers.InitGenesis(ctx, *k, genesisState)
	got := cosmoscheckers.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.NextGame, got.NextGame)
	require.ElementsMatch(t, genesisState.StoredGameList, got.StoredGameList)
	// this line is used by starport scaffolding # genesis/test/assert
}
