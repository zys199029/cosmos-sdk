package nameservice

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgSetName struct {
	Name  string
	Value string
	Owner sdk.AccAddress
}

func NewMsgSetName(name string, value string, owner sdk.AccAddress) MsgSetName {
	return MsgSetName{
		Name:  name,
		Value: value,
		Owner: owner,
	}
}

// Implements Msg.
func (msg MsgSetName) Type() string { return "nameservice" }

// Implements Msg.
func (msg MsgSetName) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("Name and Value cannot be empty")
	}
	return nil
}

// Implements Msg.
func (msg MsgSetName) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
