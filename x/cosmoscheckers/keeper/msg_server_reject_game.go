package keeper

import (
	"context"
	"strings"

	rules "github.com/colincassens/cosmosCheckers/x/cosmoscheckers/rules"
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RejectGame(goCtx context.Context, msg *types.MsgRejectGame) (*types.MsgRejectGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.IdValue)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "game not found %s", msg.IdValue)
	}

	//If game is already won you cant reject
	if storedGame.Winner != rules.NO_PLAYER.Color {
		return nil, types.ErrGameFinished
	}

	//Determine if the sender is a player and if they have already played
	if strings.Compare(storedGame.Red, msg.Creator) == 0 {
		if 1 < storedGame.MoveCount {
			return nil, types.ErrRedAlreadyPlayed
		}
	} else if strings.Compare(storedGame.Black, msg.Creator) == 0 {
		if 0 < storedGame.MoveCount {
			return nil, types.ErrBlackAlreadyPlayed
		}
	} else {
		return nil, types.ErrCreatorNotPlayer
	}

	// Remove from the FIFO
	nextGame, found := k.Keeper.GetNextGame(ctx)
	if !found {
		panic("NextGame not found")
	}
	k.Keeper.RemoveFromFifo(ctx, &storedGame, &nextGame)
	k.Keeper.MustRefundWager(ctx, &storedGame)

	//Remove the game
	k.Keeper.RemoveStoredGame(ctx, msg.IdValue)
	k.Keeper.SetNextGame(ctx, nextGame)

	//Send Event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "checkers"),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.RejectGameEventKey),
			sdk.NewAttribute(types.RejectGameEventCreator, msg.Creator),
			sdk.NewAttribute(types.RejectGameEventIdValue, msg.IdValue),
		),
	)

	//Consume Gas
	ctx.GasMeter().ConsumeGas(types.RejectGameGas, "Reject Game Gas")

	return &types.MsgRejectGameResponse{}, nil
}
