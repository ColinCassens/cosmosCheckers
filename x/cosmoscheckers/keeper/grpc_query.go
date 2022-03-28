package keeper

import (
	"github.com/colincassens/cosmosCheckers/x/cosmoscheckers/types"
)

var _ types.QueryServer = Keeper{}
