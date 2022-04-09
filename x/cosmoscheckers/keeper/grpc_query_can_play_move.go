package keeper

import (
	"context"
	"strings"
	"fmt"

	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rules "github.com/colincassens/cosmosCheckers/x/cosmoscheckers/rules"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) CanPlayMove(goCtx context.Context, req *types.QueryCanPlayMoveRequest) (*types.QueryCanPlayMoveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	storedGame, found := k.GetStoredGame(ctx, req.IdValue)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, types.ErrGameNotFound.Error(), req.IdValue)
	}

	//Make sure game is active
	if storedGame.Winner != rules.NO_PLAYER.Color {
		return &types.QueryCanPlayMoveResponse{
			Possible: false,
			Reason:  types.ErrGameFinished.Error(),
		}, nil
	}

	//Get Player
	var player rules.Player
	if strings.Compare(rules.BLACK_PLAYER.Color, req.Player) == 0 {
		player = rules.BLACK_PLAYER
	} else if strings.Compare(rules.RED_PLAYER.Color, req.Player) == 0 {
		player = rules.RED_PLAYER
	} else {
		return &types.QueryCanPlayMoveResponse{
			Possible: false,
			Reason: types.ErrCreatorNotPlayer.Error(),
		}, nil
	}

	//Get game and check player turn
	game, err := storedGame.ParseGame()
	if err != nil {
		return nil, err
	}

	if !game.TurnIs(player) {
		return &types.QueryCanPlayMoveResponse{
			Possible: false,
			Reason: types.ErrNotPlayerTurn.Error(),
		}, nil
	}

	//Check if the move is valid and return accordingly
	_, moveErr := game.Move(
		rules.Pos{
			X: int(req.FromX),
			Y: int(req.FromY),
		},
		rules.Pos{
			X: int(req.ToX),
			Y: int(req.ToY),
		},
	)

	if moveErr != nil {
		return &types.QueryCanPlayMoveResponse{
			Possible: false,
			Reason: fmt.Sprintf(types.ErrWrongMove.Error(), moveErr.Error()),
		}, nil
	}

	return &types.QueryCanPlayMoveResponse{
		Possible: true,
		Reason: "Valid Move",
	}, nil
}
