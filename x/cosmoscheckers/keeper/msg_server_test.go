package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/colincassens/cosmosCheckers/testutil/keeper"
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/keeper"
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.CosmoscheckersKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
