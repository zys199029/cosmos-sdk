package slashing

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// Expected interfaces for a StakeKeeper passed into the governance Keeper
type ValidatorSet interface {
	GetBondedValidatorOwnerAddresses(sdk.Context) []sdk.AccAddress // return a slice of owner Address of all bonded validators
	GetTotalPower(sdk.Context) sdk.Rat                             // total power of the validator set

	GetValidatorPubKey(sdk.Context, sdk.AccAddress) (crypto.PubKey, error)          // Gets the signing key of the Validator
	GetValidatorPower(sdk.Context, sdk.AccAddress) (sdk.Rat, error)                 // validation power
	GetValidatorTotalDelegationShares(sdk.Context, sdk.AccAddress) (sdk.Rat, error) // Total out standing delegator shares
	GetValidatorBondHeight(sdk.Context, sdk.AccAddress) (int64, error)              // height in which the validator became active
	ValidatorIsRevoked(sdk.Context, sdk.AccAddress) (bool, error)                   // Returns true if Validator is revoked

	GetDelegatorDelegations(sdk.Context, sdk.AccAddress) []sdk.AccAddress
	GetDelegatorDelegationShares(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.AccAddress) (sdk.Rat, error)

	// slash the validator and delegators of the validator, specifying offence height, offence power, and slash fraction
	Slash(sdk.Context, crypto.PubKey, int64, int64, sdk.Rat)
	Revoke(sdk.Context, crypto.PubKey)   // revoke a validator
	Unrevoke(sdk.Context, crypto.PubKey) // unrevoke a validator
}
