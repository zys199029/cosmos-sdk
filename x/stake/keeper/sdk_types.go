package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake/types"
)

// iterate through the active validator set and perform the provided function
func (k Keeper) GetValidatorOwnerAddresses(ctx sdk.Context) (ownerAddresses []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	for ; iterator.Valid(); iterator.Next() {
		ownerAddresses = append(ownerAddresses, sdk.AccAddress(iterator.Key()[1:]))
	}
	iterator.Close()
	return
}

// iterate through the active validator set and perform the provided function
func (k Keeper) GetBondedValidatorOwnerAddresses(ctx sdk.Context) (ownerAddresses []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsBondedIndexKey)
	for ; iterator.Valid(); iterator.Next() {
		ownerAddresses = append(ownerAddresses, GetAddressFromValBondedIndexKey(iterator.Key()))
	}
	iterator.Close()
	return
}

// get the sdk.validator for a particular address
func (k Keeper) Validator(ctx sdk.Context, address sdk.AccAddress) (types.Validator, error) {
	val, found := k.GetValidator(ctx, address)
	if !found {
		return types.Validator{}, error.Error("validator now found")
	}
	return val
}

// total power from the bond
func (k Keeper) GetTotalPower(ctx sdk.Context) sdk.Rat {
	pool := k.GetPool(ctx)
	return pool.BondedTokens
}

// get the delegation for a particular set of delegator and validator addresses
func (k Keeper) Delegation(ctx sdk.Context, addrDel sdk.AccAddress, addrVal sdk.AccAddress) (types.Delegation, error) {
	bond, ok := k.GetDelegation(ctx, addrDel, addrVal)
	if !ok {
		return types.Delegation{}, errors.New("Could not find delegation")
	}
	return bond, nil
}

// iterate through the active validator set and perform the provided function
func (k Keeper) IterateDelegations(ctx sdk.Context, delAddr sdk.AccAddress, fn func(index int64, delegation types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	key := GetDelegationsKey(delAddr)
	iterator := sdk.KVStorePrefixIterator(store, key)
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Key(), iterator.Value())
		stop := fn(i, delegation) // XXX is this safe will the fields be able to get written to?
		if stop {
			break
		}
		i++
	}
	iterator.Close()
}
