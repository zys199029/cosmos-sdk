package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

// TODO: Codespace will be removed
const DefaultCodespace = 65534

type Keeper struct {
	key sdk.StoreKey

	bk   bank.Keeper
	ibck ibc.Keeper
}

func NewKeeper(key sdk.StoreKey, bk bank.Keeper, ibck ibc.Keeper) Keeper {
	return Keeper{
		key:  key,
		bk:   bk,
		ibck: ibck,
	}
}

func (k Keeper) ibcStore(ctx sdk.Context) sdk.KVStore {
	// Prefixing for future compatability
	return ctx.KVStore(k.key).Prefix([]byte{0x00})
}
