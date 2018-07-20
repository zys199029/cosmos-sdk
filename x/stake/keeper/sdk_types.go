package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake/types"
	"github.com/tendermint/tendermint/crypto"
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

// returns whether or not a validator is revoked
func (k Keeper) ValidatorIsRevoked(ctx sdk.Context, address sdk.AccAddress) (bool, error) {
	val, found := k.GetValidator(ctx, address)
	if !found {
		return false, errors.New("validator not found")
	}
	return val.GetRevoked(), nil
}

// returns a validator's moniker
func (k Keeper) GetValidatorMoniker(ctx sdk.Context, address sdk.AccAddress) (string, error) {
	val, found := k.GetValidator(ctx, address)
	if !found {
		return "", errors.New("validator not found")
	}
	return val.GetMoniker(), nil
}

// returns a validator's status
func (k Keeper) GetValidatorStatus(ctx sdk.Context, address sdk.AccAddress) (sdk.BondStatus, error) {
	val, found := k.GetValidator(ctx, address)
	if !found {
		return 0xff, errors.New("validator not found")
	}
	return val.GetStatus(), nil
}

// returns a validator's pubkey
func (k Keeper) GetValidatorPubKey(ctx sdk.Context, address sdk.AccAddress) (crypto.PubKey, error) {
	val, found := k.GetValidator(ctx, address)
	if !found {
		return nil, errors.New("validator not found")
	}
	return val.GetPubKey(), nil
}

// returns a validator's power
func (k Keeper) GetValidatorPower(ctx sdk.Context, address sdk.AccAddress) (sdk.Rat, error) {
	val, found := k.GetValidator(ctx, address)
	if !found {
		return sdk.ZeroRat(), errors.New("validator not found")
	}
	return val.GetPower(), nil
}

// Total out standing delegator shares
func (k Keeper) GetValidatorTotalDelegationShares(ctx sdk.Context, address sdk.AccAddress) (sdk.Rat, error) {
	val, found := k.GetValidator(ctx, address)
	if !found {
		return sdk.ZeroRat(), errors.New("validator not found")
	}
	return val.GetDelegatorShares(), nil
}

// height in which the validator became active
func (k Keeper) GetValidatorBondHeight(ctx sdk.Context, address sdk.AccAddress) (int64, error) {
	val, found := k.GetValidator(ctx, address)
	if !found {
		return 0, errors.New("validator not found")
	}
	return val.GetBondHeight(), nil
}

// total voting power
func (k Keeper) GetTotalPower(ctx sdk.Context) sdk.Rat {
	pool := k.GetPool(ctx)
	return pool.BondedTokens
}

// Returns the shares that a delegator has in a certain validator pool
func (k Keeper) GetDelegatorDelegationShares(ctx sdk.Context, delegatorAddress sdk.AccAddress, validatorAddress sdk.AccAddress) (sdk.Rat, error) {
	bond, ok := k.GetDelegation(ctx, delegatorAddress, validatorAddress)
	if !ok {
		return sdk.ZeroRat(), errors.New("Could not find delegation")
	}
	return bond.GetBondShares(), nil
}

// Returns the shares that a delegator has in a certain validator pool
func (k Keeper) GetDelegatorDelegationHeight(ctx sdk.Context, delegatorAddress sdk.AccAddress, validatorAddress sdk.AccAddress) (int64, error) {
	bond, ok := k.GetDelegation(ctx, delegatorAddress, validatorAddress)
	if !ok {
		return -1, errors.New("Could not find delegation")
	}
	return bond.Height, nil
}

// Returns a slice of all the validator addresses that a certain delegator is delegated to
func (k Keeper) IterateDelegatorDelegations(ctx sdk.Context, delAddr sdk.AccAddress) (validatorAddresses []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := GetDelegationsKey(delAddr)
	iterator := sdk.KVStorePrefixIterator(store, key)
	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Key(), iterator.Value())
		validatorAddresses = append(validatorAddresses, delegation.ValidatorAddr)
	}
	iterator.Close()
	return
}
