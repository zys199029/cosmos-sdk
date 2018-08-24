package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case "egress":
			return queryEgress(ctx, path[1:], req, keeper)
		case "ingress":
			return queryIngress(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown ibc query endpoint")
		}
	}
}

type QueryEgressParams struct {
	DestChain string
	Index     int64
}

func queryEgress(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryEgressParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return []byte{}, sdk.ErrUnknownRequest("incorrectly formatted request data")
	}

}

type QueryIngressParams struct {
	SrcChain string
}

func queryIngress(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryIngressParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return []byte{}, sdk.ErrUnknownRequest("incorrectly formatted request data")
	}
}
