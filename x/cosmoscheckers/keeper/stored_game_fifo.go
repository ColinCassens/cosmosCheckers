package keeper

import (
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) RemoveFromFifo(ctx sdk.Context, game *types.StoredGame, info *types.NextGame) {
	if game.BeforeId != types.NoFifoIdKey {
		beforeElement, found := k.GetStoredGame(ctx, game.BeforeId)
		if !found {
			panic("Prev Game Not Found! Linked List Error")
		}

		beforeElement.AfterId = game.AfterId
		k.SetStoredGame(ctx, beforeElement)

		if game.AfterId == types.NoFifoIdKey {
			info.FifoTail = beforeElement.Index
		}
	}
	if game.AfterId != types.NoFifoIdKey {
		afterElement, found := k.GetStoredGame(ctx, game.AfterId)
		if !found {
			panic("Following Game Not Found! Linked List Error")
		}
		afterElement.BeforeId = game.BeforeId
		k.SetStoredGame(ctx, afterElement)

		if game.BeforeId == types.NoFifoIdKey {
			info.FifoHead = afterElement.Index
		}
	}
}

func (k Keeper) SendToFifoTail(ctx sdk.Context, game *types.StoredGame, info *types.NextGame) {
	//Check for empty LL and act accordingly
	if info.FifoTail == types.NoFifoIdKey && info.FifoHead == types.NoFifoIdKey {
		info.FifoHead = game.Index
		info.FifoTail = game.Index
		game.BeforeId = types.NoFifoIdKey
		game.AfterId = types.NoFifoIdKey
	} else if info.FifoTail == types.NoFifoIdKey || info.FifoHead == types.NoFifoIdKey {
		//Check for head or tail but not the other (ERROR)
		panic("Linked List Error!! Fifo needs both head & tail")
	} else if !(info.FifoTail == game.Index) {
		// Make sure Game isnt already tail ^
		//Remove from LL
		k.RemoveFromFifo(ctx, game, info)

		//Add to the tail of the fifo
		tailElement, found := k.GetStoredGame(ctx, info.FifoTail)
		if !found {
			panic("Tail of LL not found!")
		}

		tailElement.AfterId = game.Index
		game.BeforeId = tailElement.Index
		info.FifoTail = game.Index
	}
}
