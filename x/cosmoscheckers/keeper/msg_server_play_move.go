package keeper

import (
	"context"
	"strconv"
	"strings"

	rules "github.com/colincassens/cosmosCheckers/x/cosmoscheckers/rules"
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) PlayMove(goCtx context.Context, msg *types.MsgPlayMove) (*types.MsgPlayMoveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//Fetch Stored game
	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.IdValue)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "game not found %s", msg.IdValue)
	}

	//Check player legitimacy
	var player rules.Player
	if strings.Compare(storedGame.Red, msg.Creator) == 0 {
		player = rules.RED_PLAYER
	} else if strings.Compare(storedGame.Black, msg.Creator) == 0 {
		player = rules.BLACK_PLAYER
	} else {
		return nil, types.ErrCreatorNotPlayer
	}

	//Instantiate the game
	game, err := storedGame.ParseGame()
	if err != nil {
		panic(err.Error())
	}

	//Check Player turn
	if !game.TurnIs(player) {
		return nil, types.ErrNotPlayerTurn
	}

	//Move
	captured, moveErr := game.Move(
		rules.Pos{
			X: int(msg.FromX),
			Y: int(msg.FromY),
		},
		rules.Pos{
			X: int(msg.ToX),
			Y: int(msg.ToY),
		},
	)
	if moveErr != nil {
		return nil, sdkerrors.Wrapf(moveErr, types.ErrWrongMove.Error())
	}

	//Increment Move Counter
	storedGame.MoveCount++
	storedGame.Deadline = types.FormatDeadline(types.GetNextDeadline(ctx))

	// Send to the back of the FIFO
	nextGame, found := k.Keeper.GetNextGame(ctx)
	if !found {
		panic("NextGame not found")
	}
	k.Keeper.SendToFifoTail(ctx, &storedGame, &nextGame)

	//Store Info
	storedGame.Game = game.String()
	storedGame.Turn = game.Turn.Color
	k.Keeper.SetStoredGame(ctx, storedGame)
	k.Keeper.SetNextGame(ctx, nextGame)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "checkers"),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.PlayMoveEventKey),
			sdk.NewAttribute(types.PlayMoveEventCreator, msg.Creator),
			sdk.NewAttribute(types.PlayMoveEventIdValue, msg.IdValue),
			sdk.NewAttribute(types.PlayMoveEventCapturedX, strconv.FormatInt(int64(captured.X), 10)),
			sdk.NewAttribute(types.PlayMoveEventCapturedY, strconv.FormatInt(int64(captured.Y), 10)),
			sdk.NewAttribute(types.PlayMoveEventWinner, game.Winner().Color),
		),
	)

	//Return Relevant info
	return &types.MsgPlayMoveResponse{
		IdValue:   msg.IdValue,
		CapturedX: int64(captured.X),
		CapturedY: int64(captured.Y),
		Winner:    game.Winner().Color,
	}, nil
}
