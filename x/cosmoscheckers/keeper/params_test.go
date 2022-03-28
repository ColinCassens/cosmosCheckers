package keeper_test

import (
	"testing"

	testkeeper "github.com/colincassens/cosmosCheckers/testutil/keeper"
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.CosmoscheckersKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
