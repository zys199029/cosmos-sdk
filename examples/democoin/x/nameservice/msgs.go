package nameservice

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgBuyName struct {
	Name  string
	Value string
	Owner sdk.AccAddress
}

func NewMsgBuyName(name string, value string, owner sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		Name:  name,
		Value: value,
		Owner: owner,
	}
}

// Implements Msg.
func (msg MsgBuyName) Type() string { return "nameservice" }

// Implements Msg.
func (msg MsgBuyName) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("Name and Value cannot be empty")
	}
	return nil
}

// Implements Msg.
func (msg MsgBuyName) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
