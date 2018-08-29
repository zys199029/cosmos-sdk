package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (k Keeper) Query(ctx sdk.Context, store sdk.KVStore, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
	switch path[0] {
	case "egress-datagram":
		return queryEgressDatagram(ctx, store, path[1:], req, k)
	case "egress-index":
		return queryEgressSequence(ctx, store, path[1:], req, k)
	case "ingress-index":
		return queryIngressSequence(ctx, store, path[1:], req, k)
	default:
		return nil, sdk.ErrUnknownRequest("unknown ibc query endpoint")
	}
}

type QueryEgressDatagramParams struct {
	DestChain    string
	DatagramType byte
	Sequence     uint64
}

func queryEgressDatagram(ctx sdk.Context, store sdk.KVStore, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryEgressDatagramParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return []byte{}, sdk.ErrUnknownRequest("incorrectly formatted request data")
	}

	ty := DatagramType(params.DatagramType)
	if !ty.IsValid() {
		return []byte{}, sdk.ErrUnknownRequest("invalid datagram type")
	}

	r := keeper.channelRuntime(ctx, store, ty, params.DestChain)
	data := r.getEgressDatagram(params.Sequence)

	bz, err2 := keeper.cdc.MarshalJSON(data)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

type QueryEgressSequenceParams struct {
	DestChain    string
	DatagramType byte
}

func queryEgressSequence(ctx sdk.Context, store sdk.KVStore, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryEgressSequenceParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return []byte{}, sdk.ErrUnknownRequest("incorrectly formatted request data")
	}

	ty := DatagramType(params.DatagramType)
	if !ty.IsValid() {
		return []byte{}, sdk.ErrUnknownRequest("invalid datagram type")
	}

	r := keeper.channelRuntime(ctx, store, ty, params.DestChain)
	seq := r.getEgressSequence()

	bz, err2 := keeper.cdc.MarshalJSON(seq)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

type QueryIngressSequenceParams struct {
	SrcChain     string
	DatagramType byte
}

func queryIngressSequence(ctx sdk.Context, store sdk.KVStore, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryIngressSequenceParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return []byte{}, sdk.ErrUnknownRequest("incorrectly formatted request data")
	}

	ty := DatagramType(params.DatagramType)
	if !ty.IsValid() {
		return []byte{}, sdk.ErrUnknownRequest("invalid datagram type")
	}

	r := keeper.channelRuntime(ctx, store, ty, params.SrcChain)
	index := r.getIngressSequence()

	bz, err2 := keeper.cdc.MarshalJSON(index)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}
