package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper - handlers sets/gets of custom variables for your module
type Keeper struct {
	ck bank.Keeper

	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.

	codespace sdk.CodespaceType
}

// NewKeeper - Returns the Keeper
func NewKeeper(key sdk.StoreKey, bankKeeper bank.Keeper, codespace sdk.CodespaceType) Keeper {
	return Keeper{bankKeeper, key, codespace}
}

// GetTrend - returns the current cool trend
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(name))
	return string(bz)
}

// GetTrend - returns the current cool trend
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), []byte(value))
}

// Implements sdk.AccountMapper.
func (k Keeper) IteratePrefix(ctx sdk.Context, namePrefix string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(namePrefix))
}
