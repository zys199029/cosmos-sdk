package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/lib"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

// TODO: Codespace will be removed
const DefaultCodespace = 65534

type Keeper struct {
	key sdk.StoreKey
	cdc *wire.Codec

	bk   bank.Keeper
	ibck ibc.Keeper
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, bk bank.Keeper, ibck ibc.Keeper) Keeper {
	return Keeper{
		key:  key,
		cdc:  cdc,
		bk:   bk,
		ibck: ibck,
	}
}

// ----------------------------------
// Store Accessors

func IBCStorePrefix() []byte {
	return []byte{0x00}
}

func (k Keeper) ibcStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key).Prefix(IBCStorePrefix())
}

func LedgerPrefix() []byte {
	return []byte{0x01}
}

func (k Keeper) ledger(store sdk.KVStore) lib.Mapping {
	return lib.NewMapping(k.cdc, store.Prefix(LedgerPrefix()))
}
