package keeper

import (
	"context"
	"strconv"

	rules "github.com/colincassens/cosmosCheckers/x/cosmoscheckers/rules"
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	nextGame, found := k.Keeper.GetNextGame(ctx)
	if !found {
		panic("NextGame not found")
	}
	newIndex := strconv.FormatUint(nextGame.IdValue, 10)
	storedGame := types.StoredGame{
		Creator:   	msg.Creator,
		Index:     	newIndex,
		Game:      	rules.New().String(),
		Red:       	msg.Red,
		Black:     	msg.Black,
		MoveCount: 	0,
		Deadline:  	types.FormatDeadline(types.GetNextDeadline(ctx)),
		Winner:    	rules.NO_PLAYER.Color,
		Wager:     	msg.Wager,
		Token:	   	msg.Token,
	}

	err := storedGame.Validate()
	if err != nil {
		return nil, err
	}
	k.Keeper.SendToFifoTail(ctx, &storedGame, &nextGame)
	k.Keeper.SetStoredGame(ctx, storedGame)

	nextGame.IdValue++
	k.Keeper.SetNextGame(ctx, nextGame)

	//Signal a game created event for users to listen for
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "cosmosCheckers"),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.StoredGameEventKey),
			sdk.NewAttribute(types.StoredGameEventCreator, msg.Creator),
			sdk.NewAttribute(types.StoredGameEventIndex, newIndex),
			sdk.NewAttribute(types.StoredGameEventRed, msg.Red),
			sdk.NewAttribute(types.StoredGameEventBlack, msg.Black),
			sdk.NewAttribute(types.StoredGameEventWager, strconv.FormatUint(msg.Wager, 10)),
			sdk.NewAttribute(types.StoredGameEventToken, msg.Token),
		),
	)

	//Consume Gas
	ctx.GasMeter().ConsumeGas(types.CreateGameGas, "Create Game Gas")

	return &types.MsgCreateGameResponse{
		IdValue: newIndex,
	}, nil
}
