# Cosmos SDK Code Run #1

Code don't lie, check out why the Cosmos SDK is the ideal platform for scalable blockchain development.

## KVStore

The KVStore provides the basic persistence layer for your SDK application.

https://github.com/cosmos/cosmos-sdk/blob/3fc7200f1d1045a19efc30395e5916f9ef1b42b7/types/store.go#L91-L121

You can mount multiple KVStores onto your application, e.g. one for staking, one for accounts, one for IBC, and so on.

https://github.com/cosmos/cosmos-sdk/blob/3fc7200f1d1045a19efc30395e5916f9ef1b42b7/examples/basecoin/app/app.go#L90

The implementation of a KVStore is responsible for providing any Merkle proofs for each query, if requested.

https://github.com/cosmos/cosmos-sdk/blob/3fc7200f1d1045a19efc30395e5916f9ef1b42b7/store/iavlstore.go#L135

Stores can be cache-wrapped to provide transactions at the persistence level
(and this is well supported for iterators as well). This feature is used to
provide a layer of transactional isolation for transaction processing after the
"AnteHandler" deducts any associated fees for the transaction.  Cache-wrapping
can also be useful when implementing a virtual-machine or scripting environment
for the blockchain.

## go-amino

The Cosmos-SDK uses
[go-amino](https://github.com/cosmos/cosmos-sdk/blob/96451b55fff107511a65bf930b81fb12bed133a1/examples/basecoin/app/app.go#L97-L111)
extensively to serialize and deserialize Go types into Protobuf3 compatible
bytes.

Go-amino (e.g. over https://github.com/golang/protobuf) uses reflection to
encode/decode any Go object.  This lets the SDK developer focus on defining
data structures in Go without the need to maintain a separate schema for
Proto3.  In addition, Amino extends Proto3 with native support for interfaces
and concrete types.

For example, the Cosmos SDK's `x/auth` package imports the PubKey interface
from `tendermint/go-crypto` , where PubKey implementations include those for
Ed25519 and Secp256k1.  Each auth.BaseAccount has a PubKey.

https://github.com/cosmos/cosmos-sdk/blob/d309abb4b9bdb01272e54e048063502110c801fa/x/auth/account.go#L35-L43

Amino knows what concrete type to decode for each interface value
based on what concretes are registered for the interface.

For example, the "Basecoin" example app knows about Ed25519 and Secp256k1 keys
because they are registered by the app's codec below:

https://github.com/cosmos/cosmos-sdk/blob/d309abb4b9bdb01272e54e048063502110c801fa/examples/basecoin/app/app.go#L102-L106

For more information on Go-Amino, see https://github.com/tendermint/go-amino.

## Keys, Keepers, and Mappers

The Cosmos SDK is designed to enable an ecosystem of libraries that can be
imported together to form a whole application.  To make this ecosystem
more secure, we've developed a design pattern that follows the principle of 
least-authority.

Mappers and Keepers provide access to KV stores via the context.  The only
difference between the two is that a Mapper provides a lower-level API, so
generally a Keeper might hold references to other Keepers and Mappers but not
vice versa.

Mappers and Keepers don't hold any references to any stores directly.  They only
hold a "key" (the `sdk.StoreKey` below):

https://github.com/cosmos/cosmos-sdk/blob/fc0e4013278d41fab4f3ac73f28a42bc45889106/x/auth/mapper.go#L14-L24

This way, you can hook everything up in your main app.go file and see what
components have access to what stores and other components.

https://github.com/cosmos/cosmos-sdk/blob/d309abb4b9bdb01272e54e048063502110c801fa/examples/basecoin/app/app.go#L65-L70

Later during the execution of a transaction (e.g. via ABCI DeliverTx after a
block commit) the context is passed in as the first argument.  The context
includes references to any relevant KV stores, but you can only access them if
you hold the associated key.

https://github.com/cosmos/cosmos-sdk/blob/fc0e4013278d41fab4f3ac73f28a42bc45889106/x/auth/mapper.go#L44-L53

Mappers and Keepers cannot hold direct references to stores because the store
is not known at app initialization time.  The store is dynamically created (and
wrapped via CacheKVStore as needed to provide a transactional context) for
every committed transaction (via ABCI DeliverTx) and mempool check transaction
(via ABCI CheckTx). 

## Tx, Msg, Handler, and AnteHandler

A transaction (Tx interface) is a signed/authenticated message (Msg interface).

Transactions that are discovered by the Tendermint mempool are processed by the
AnteHandler ("ante" just means "before") where the validity of the transaction
is checked and any fees are collected.

Transactions that get committed in a block first get processed through the
AnteHandler, and if the transaction is valid after fees are deducted, they are
processed through the appropriate Handler.

In either case, the transaction bytes must first be parsed.  The default
transaction parser uses Amino.  Most SDK developers will want to use the
standard transaction structure defined in the `x/auth` package (and the
corresponding AnteHandler implementation also provided in `x/auth`):

https://github.com/cosmos/cosmos-sdk/blob/fc0e4013278d41fab4f3ac73f28a42bc45889106/x/auth/stdtx.go#L12-L18

Various packages generally define their own message types.  The Basecoin
example app includes multiple message types that are registered in app.go:

https://github.com/cosmos/cosmos-sdk/blob/d309abb4b9bdb01272e54e048063502110c801fa/examples/basecoin/app/app.go#L102-L106

Finally, handlers are added to the router in your app.go file to map messages
to their corresponding handlers. (In the future we will provide more routing
features to enable pattern matching for more flexibility).

https://github.com/cosmos/cosmos-sdk/blob/d309abb4b9bdb01272e54e048063502110c801fa/examples/basecoin/app/app.go#L78-L83

## EndBlocker

The EndBlocker hook allows us to register callback logic to be performed at the
end of each block.  This lets us process background events, such as processing
validator inflationary atom provisions:

https://github.com/cosmos/cosmos-sdk/blob/3fc7200f1d1045a19efc30395e5916f9ef1b42b7/x/stake/handler.go#L32-L37

By the way, the SDK provides a staking module, which provides all the
bonding/unbonding funcionality for the Cosmos Hub:
https://github.com/cosmos/cosmos-sdk/tree/develop/x/stake (staking module)

## See Also:

https://tendermint.com
https://cosmos.network
https://github.com/cosmos/cosmos-sdk

We're hiring!
