package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/cosmos/cosmos-sdk/x/ibc"

	"github.com/cosmos/cosmos-sdk/x/bank/ibc"
)

type transferBody struct {
	// Fees             sdk.Coin  `json="fees"`
	Amount           sdk.Coins `json:"amount"`
	LocalAccountName string    `json:"name"`
	Password         string    `json:"password"`
	SrcChainID       string    `json:"src_chain_id"`
	AccountNumber    int64     `json:"account_number"`
	Sequence         int64     `json:"sequence"`
	Gas              int64     `json:"gas"`
}

// TransferRequestHandler - http request handler to transfer coins to a address
// on a different chain via IBC
func TransferRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		destChainID := vars["destchain"]
		bech32addr := vars["address"]

		to, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(&w, http.StatusBadRequest, err.Error())
			return
		}

		var m transferBody
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.WriteErrorResponse(&w, http.StatusBadRequest, err.Error())
			return
		}

		err = cdc.UnmarshalJSON(body, &m)
		if err != nil {
			utils.WriteErrorResponse(&w, http.StatusBadRequest, err.Error())
			return
		}

		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			utils.WriteErrorResponse(&w, http.StatusUnauthorized, err.Error())
			return
		}

		from, err := sdk.AccAddressFromBech32(string(info.GetPubKey().Address()))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// build message
		p := bank.PayloadCoins{
			SrcAddr:  from,
			DestAddr: to,
			Coins:    m.Amount,
		}

		msg := ibc.MsgSend{
			Payload:   p,
			DestChain: destChainID,
		}

		txCtx := authctx.TxContext{
			Codec:         cdc,
			ChainID:       m.SrcChainID,
			AccountNumber: m.AccountNumber,
			Sequence:      m.Sequence,
			Gas:           m.Gas,
		}

		if m.Gas == 0 {
			newCtx, err := utils.EnrichCtxWithGas(txCtx, cliCtx, m.LocalAccountName, m.Password, []sdk.Msg{msg})
			if err != nil {
				utils.WriteErrorResponse(&w, http.StatusInternalServerError, err.Error())
				return
			}
			txCtx = newCtx
		}

		txBytes, err := txCtx.BuildAndSign(m.LocalAccountName, m.Password, []sdk.Msg{msg})
		if err != nil {
			utils.WriteErrorResponse(&w, http.StatusUnauthorized, err.Error())
			return
		}

		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			utils.WriteErrorResponse(&w, http.StatusInternalServerError, err.Error())
			return
		}

		output, err := cdc.MarshalJSON(res)
		if err != nil {
			utils.WriteErrorResponse(&w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}
