package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/lib"
	"github.com/cosmos/cosmos-sdk/wire"
)

// ------------------------------------------
// Type Definitions

type DatagramType byte

const (
	PacketType = DatagramType(iota)
	ReceiptType
	DatagramTypeLength
)

func (ty DatagramType) IsValid() bool {
	return ty < DatagramTypeLength
}

type Header struct {
	SrcChain  string
	DestChain string
}

func (h Header) InverseDirection() Header {
	return Header{
		SrcChain:  h.DestChain,
		DestChain: h.SrcChain,
	}
}

type Payload interface {
	Type() string
	ValidateBasic() sdk.Error
	GetSigners() []sdk.AccAddress
	DatagramType() DatagramType
}

type Datagram struct {
	Header
	// Should we unexport Payload to possible modification from the modules?
	Payload
}

type Proof struct {
	Height   uint64
	Sequence uint64
}

// -------------------------------------------
// Store Accessors

func EgressQueuePrefix(ty DatagramType, chainid string) []byte {
	return append(append([]byte{0x00}, byte(ty)), []byte(chainid)...)
}

func egressQueue(store sdk.KVStore, cdc *wire.Codec, ty DatagramType, chainid string) lib.Linear {
	return lib.NewLinear(cdc, store.Prefix(EgressQueuePrefix(ty, chainid)), nil)
}

func IngressSequenceKey(ty DatagramType, chainid string) []byte {
	return append(append([]byte{0x01}, byte(ty)), []byte(chainid)...)
}

func ingressSequence(store sdk.KVStore, cdc *wire.Codec, ty DatagramType, chainid string) lib.Value {
	return lib.NewValue(store, cdc, IngressSequenceKey(ty, chainid))
}

// --------------------------------------------
// Channel Runtime

type channelRuntime struct {
	k               Keeper
	egressQueue     lib.List
	ingressSequence lib.Value
	thisChain       string
	thatChain       string
}

func (k Keeper) channelRuntime(ctx sdk.Context, store sdk.KVStore, ty DatagramType, thatChain string) channelRuntime {
	return channelRuntime{
		k:               k,
		egressQueue:     egressQueue(store, k.cdc, ty, thatChain),
		ingressSequence: ingressSequence(store, k.cdc, ty, thatChain),
		thisChain:       ctx.ChainID(),
		thatChain:       thatChain,
	}
}

func (r channelRuntime) getEgressSequence() uint64 {
	return r.egressQueue.Len()
}

func (r channelRuntime) getEgressDatagram(index uint64) (res Datagram) {
	r.egressQueue.Get(index, &res)
	return
}

func (r channelRuntime) pushEgressDatagram(data Datagram) {
	r.egressQueue.Push(data)
}

func (r channelRuntime) getIngressSequence() (res uint64) {
	r.ingressSequence.GetIfExists(&res)
	return
}

func (r channelRuntime) setIngressSequence(seq uint64) {
	r.ingressSequence.Set(seq)
}
