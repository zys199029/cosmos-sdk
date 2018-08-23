package nameservice

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "cool" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgQuiz This is the engine of your module
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg MsgSetName) sdk.Result {
	keeper.SetName(ctx, msg.Name, msg.Value)
	return sdk.Result{}
}
