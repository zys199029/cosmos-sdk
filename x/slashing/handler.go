package slashing

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case MsgUnrevoke:
			return handleMsgUnrevoke(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("invalid message parse in staking module").Result()
		}
	}
}

// Validators must submit a transaction to unrevoke itself after
// having been revoked (and thus unbonded) for downtime
func handleMsgUnrevoke(ctx sdk.Context, msg MsgUnrevoke, k Keeper) sdk.Result {

	isRevoked, err := k.validatorSet.ValidatorIsRevoked(ctx, msg.ValidatorAddr)
	if err != nil {
		return ErrNoValidatorForAddress(k.codespace).Result()
	}
	if !isRevoked {
		return ErrValidatorNotRevoked(k.codespace).Result()
	}

	valPubKey, _ := k.validatorSet.GetValidatorPubKey(ctx, msg.ValidatorAddr)
	addr := sdk.ValAddress(valPubKey.Address())

	// Signing info must exist
	info, found := k.getValidatorSigningInfo(ctx, addr)
	if !found {
		return ErrNoValidatorForAddress(k.codespace).Result()
	}

	// Cannot be unrevoked until out of jail
	if ctx.BlockHeader().Time < info.JailedUntil {
		return ErrValidatorJailed(k.codespace).Result()
	}

	if ctx.IsCheckTx() {
		return sdk.Result{}
	}

	// Update the starting height (so the validator can't be immediately revoked again)
	info.StartHeight = ctx.BlockHeight()
	k.setValidatorSigningInfo(ctx, addr, info)

	// Unrevoke the validator
	k.validatorSet.Unrevoke(ctx, valPubKey)

	tags := sdk.NewTags("action", []byte("unrevoke"), "validator", []byte(msg.ValidatorAddr.String()))

	return sdk.Result{
		Tags: tags,
	}
}
