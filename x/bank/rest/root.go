package rest

import (
	"github.com/gorilla/mux"

	keys "github.com/tendermint/go-crypto/keys"

	"github.com/cosmos/cosmos-sdk/wire"
)

func RegisterRoutes(r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc("/accounts/{address}/send", SendRequestHandler(cdc, kb)).Methods("POST")
}
