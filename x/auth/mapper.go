package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
)

var _ sdk.AccountMapper = (*accountMapper)(nil)
var _ sdk.AccountMapper = (*sealedAccountMapper)(nil)

// AccountConstructor returns an empty sdk.Account concrete type
type AccountConstructor = func() sdk.Account

// Implements sdk.AccountMapper.
// This AccountMapper encodes/decodes accounts using the
// go-amino (binary) encoding/decoding library.
type accountMapper struct {

	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// Constructor for the sdk.Account concrete type.
	accConstructor AccountConstructor

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec

	// The key under which the next Account Number will be stored
	nextAccNumKey []byte
}

// NewAccountMapper returns a new sdk.AccountMapper that
// uses go-amino to (binary) encode and decode concrete sdk.Accounts.
// nolint
func NewAccountMapper(cdc *wire.Codec, key sdk.StoreKey, accConstructor AccountConstructor) accountMapper {
	nextAccNumKey, _ := cdc.MarshalBinary("NextAccountNum")
	return accountMapper{
		key:            key,
		accConstructor: accConstructor,
		cdc:            cdc,
		nextAccNumKey:  nextAccNumKey,
	}
}

// Returns the go-amino codec.  You may need to register interfaces
// and concrete types here, if your app's sdk.Account
// implementation includes interface fields.
// NOTE: It is not secure to expose the codec, so check out
// .Seal().
func (am accountMapper) WireCodec() *wire.Codec {
	return am.cdc
}

// Returns a "sealed" accountMapper.
// The codec is not accessible from a sealedAccountMapper.
func (am accountMapper) Seal() sealedAccountMapper {
	return sealedAccountMapper{am}
}

// Implements sdk.AccountMapper.
func (am accountMapper) NewAccountWithAddress(ctx sdk.Context, addr sdk.Address) sdk.Account {
	acc := am.accConstructor()
	acc.SetAddress(addr)
	return acc
}

// Implements sdk.AccountMapper.
func (am accountMapper) GetAccount(ctx sdk.Context, addr sdk.Address) sdk.Account {
	store := ctx.KVStore(am.key)
	bz := store.Get(addr)
	if bz == nil {
		return nil
	}
	acc := am.decodeAccount(bz)
	return acc
}

// Implements sdk.AccountMapper.
func (am accountMapper) SetAccount(ctx sdk.Context, acc sdk.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(am.key)
	bz := am.encodeAccount(acc)
	store.Set(addr, bz)
}

//----------------------------------------
// sealedAccountMapper

type sealedAccountMapper struct {
	accountMapper
}

// There's no way for external modules to mutate the
// sam.accountMapper.cdc from here, even with reflection.
func (sam sealedAccountMapper) WireCodec() *wire.Codec {
	panic("accountMapper is sealed")
}

//----------------------------------------
// misc.

func (am accountMapper) encodeAccount(acc sdk.Account) []byte {
	bz, err := am.cdc.MarshalBinaryBare(acc)
	if err != nil {
		panic(err)
	}
	return bz
}

func (am accountMapper) decodeAccount(bz []byte) (acc sdk.Account) {
	err := am.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
}
