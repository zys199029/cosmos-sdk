package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// Expected interfaces for a StakeKeeper passed into the governance Keeper

// validator for a delegated proof of stake system
type Validator interface {
	GetOwner() sdk.AccAddress    // owner AccAddress to receive/return validators coins
	GetPower() sdk.Rat           // validation power
	GetDelegatorShares() sdk.Rat // Total out standing delegator shares
	GetBondHeight() int64        // height in which the validator became active
}

// properties for the set of all validators
type ValidatorSet interface {
	// iterate through bonded validator by pubkey-AccAddress, execute func for each validator
	IterateValidatorsBonded(sdk.Context,
		func(index int64, validator Validator) (stop bool))

	Validator(sdk.Context, sdk.AccAddress) Validator // get a particular validator by owner AccAddress
	TotalPower(sdk.Context) sdk.Rat                  // total power of the validator set

	// slash the validator and delegators of the validator, specifying offence height, offence power, and slash fraction
	Slash(sdk.Context, crypto.PubKey, int64, int64, sdk.Rat)
	Revoke(sdk.Context, crypto.PubKey)   // revoke a validator
	Unrevoke(sdk.Context, crypto.PubKey) // unrevoke a validator
}

//_______________________________________________________________________________

// delegation bond for a delegated proof of stake system
type Delegation interface {
	GetDelegator() sdk.AccAddress // delegator AccAddress for the bond
	GetValidator() sdk.AccAddress // validator owner AccAddress for the bond
	GetBondShares() sdk.Rat       // amount of validator's shares
}

// properties for the set of all delegations for a particular
type DelegationSet interface {
	GetValidatorSet() ValidatorSet // validator set for which delegation set is based upon

	// iterate through all delegations from one delegator by validator-AccAddress,
	//   execute func for each validator
	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation Delegation) (stop bool))
}
