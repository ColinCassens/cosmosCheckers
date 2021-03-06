package keeper

import (
	"fmt"
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	rules "github.com/colincassens/cosmosCheckers/x/cosmoscheckers/rules"
)

func (k *Keeper) CollectWager(ctx sdk.Context, storedGame *types.StoredGame) error {
	if storedGame.MoveCount == 0 {
		// Black pays
		black, err := storedGame.GetBlackAddress()
		if err != nil {
			panic(err.Error())
		}
		err = k.bank.SendCoinsFromAccountToModule(ctx, black, types.ModuleName, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			panic(err.Error())
		}

	} else if storedGame.MoveCount == 1 {
		// Red pays
		red, err := storedGame.GetRedAddress()
		if err != nil {
			panic(err.Error())
		}
		err = k.bank.SendCoinsFromAccountToModule(ctx, red, types.ModuleName, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			panic(err.Error())
		}

	}
	return nil
}

func (k *Keeper) MustPayWinnings(ctx sdk.Context, storedGame *types.StoredGame) {
	winner, found, err := storedGame.GetWinnerAddress()
	if err != nil {
		panic(err.Error())
	}
	if !found {
		panic(fmt.Sprintf(types.ErrCannotFindWinnerByColor.Error(), storedGame.Winner))
	}

	winnings := storedGame.GetWagerCoin()
	if storedGame.MoveCount == 0 {
		panic(types.ErrNothingToPay.Error())
	} else if 1 < storedGame.MoveCount {
		winnings = winnings.Add(winnings)
	}

	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, winner, sdk.NewCoins(winnings))
	if err != nil {
		panic(err.Error())
	}
}

func (k *Keeper) MustRefundWager(ctx sdk.Context, storedGame *types.StoredGame) {
	if storedGame.MoveCount == 1 {
		black, err := storedGame.GetBlackAddress()
		if err != nil {
			panic(err.Error())
		}
		err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, black, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			panic(fmt.Sprintf(types.ErrCannotRefundWager.Error(), rules.BLACK_PLAYER.Color))
		}

	} else if storedGame.MoveCount == 0 {
		//Do Nothing
	} else {
		//TODO: Implement Draw
		panic(fmt.Sprintf(types.ErrNotInRefundState.Error(), storedGame.MoveCount))
	}
}
