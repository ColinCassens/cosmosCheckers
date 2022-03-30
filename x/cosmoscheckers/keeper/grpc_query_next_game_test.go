package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/colincassens/cosmosCheckers/testutil/keeper"
	"github.com/colincassens/cosmosCheckers/testutil/nullify"
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/types"
)

func TestNextGameQuery(t *testing.T) {
	keeper, ctx := keepertest.CosmoscheckersKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	item := createTestNextGame(keeper, ctx)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetNextGameRequest
		response *types.QueryGetNextGameResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetNextGameRequest{},
			response: &types.QueryGetNextGameResponse{NextGame: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.NextGame(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
